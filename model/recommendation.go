package model

import (
	"math"
	"sort"
	"fmt"
	"log"
)

type Recommendation struct {
	UuNeighbours uint
}

// returns: movieID => averageScore
func (r *Recommendation) ComputeFilmAverageScores(ratings []Rating) map[int64]float64 {
	userRatings := r.getUserRatingVectors(ratings)
	filmRatings := r.getMovieUserRatings(ratings)

	userAverages := make(map[int64]float64)
	for userID, userRats := range userRatings {
		userAverages[userID] = r.MeanRating(userRats)
	}

	filmScores := make(map[int64]float64)
	for filmID, userFilmRatings := range filmRatings {
		if len(userFilmRatings) < 10 {
			filmScores[filmID] = 0
			continue
		}

		filmScores[filmID] = r.filmAverageRating(userFilmRatings, userAverages)
	}

	return filmScores
}

// returns: movieID => predictedScore
// are returned only movies, which user hasn't rated
func (r *Recommendation) PredictUserScores(userID int64, ratings []Rating) []MoviePrediction {
	userMovieRatings := r.GetUserMovieRatings(ratings)
	userAverageScore := float64(0)

	// user hasn't rated any movies, we don't predict anything
	if _, ok := userMovieRatings[userID]; !ok {
		return nil
	}

	for _, rating := range userMovieRatings[userID] {
		userAverageScore += rating
	}

	userAverageScore /= float64(len(userMovieRatings[userID]))


	r.normalizeUserMovieRatings(userMovieRatings)

	//debug(userMovieRatings, true)

	// compute cosine similarity between user and another users
	cosineSimilarities := make([]Similarity, 0, len(userMovieRatings)-1)
	for uID, uMovieRatings := range userMovieRatings {

		// will not compute similarity for itself
		if uID == userID {
			continue
		}

		cosineSimilarities = append(cosineSimilarities, Similarity{
			ID:    uID,
			Value: r.cosineSimilarity(uMovieRatings, userMovieRatings[userID]),
		})
	}

	sort.Sort(BySimilarityDesc(cosineSimilarities))

	// iterate through films which user hasn't scored and compute prediction for them
	// are used 'uuNeighbours' nearest neighbours, who rated this film
	//fmt.Printf("Cosine similarities: %d\n", len(cosineSimilarities))
	//fmt.Println(cosineSimilarities)

	filmPredictions := make([]MoviePrediction, 0, 2)
	movieUserRatings := r.getMovieUserRatings(ratings)

	var score float64
	for movieID, uRatings := range movieUserRatings {
		// if user has rated this movie, we continue
		if _, ok := uRatings[userID]; ok {
			continue
		}

		// take into consideration only top n neighbours
		neighbours := uint(0)
		nominator := float64(0)
		denominator := float64(0)
		for _, similarity := range cosineSimilarities {

			if neighbours == r.UuNeighbours {
				score = nominator / denominator + userAverageScore
				if math.IsNaN(score) {
					fmt.Println("NaN score for movie %v, user %v\n", movieID, userID)
				}
				if score > 8 {
					fmt.Printf("Movie %v, User %v\n, Score %v, denom %v, nomin %v\n", movieID, userID, score, denominator, nominator)
				}

				filmPredictions = append(filmPredictions, MoviePrediction{MovieID: movieID, PredictedScore: score})
				break
			}

			if _, ok := uRatings[similarity.ID]; ok {
				neighbours++;
				// if positive similarities are not enough
				// we don't recommend this film at all
				if similarity.Value <= 0 {
					break
				}

				denominator += similarity.Value
				nominator += similarity.Value * userMovieRatings[similarity.ID][movieID]
			}
		}

	}

	sort.Sort(ByScoreDesc(filmPredictions))

	return filmPredictions
}

// return: userID => [movieRatingi, movieRating]
func (r *Recommendation) getUserRatingVectors(ratings []Rating) map[int64][]float64 {
	userRatings := make(map[int64][]float64, 0)
	for _, rating := range ratings {
		if _, ok := userRatings[rating.UserID]; ok {
			userRatings[rating.UserID] = append(userRatings[rating.UserID], rating.Value)
		} else {
			userRatings[rating.UserID] = make([]float64, 1)
			userRatings[rating.UserID][0] = rating.Value
		}
	}

	return userRatings
}

// returns: movieID => [userID1 => rating, userID2 => rating]
func (r *Recommendation) getMovieUserRatings(ratings []Rating) map[int64]map[int64]float64 {
	filmRatings := make(map[int64]map[int64]float64, 0)
	for _, rating := range ratings {
		if _, ok := filmRatings[rating.MovieID]; !ok {
			filmRatings[rating.MovieID] = make(map[int64]float64)
		}

		filmRatings[rating.MovieID][rating.UserID] = rating.Value
	}

	return filmRatings
}

// returns: userID => [movieID1 => rating, movieID2 => rating]
func (r *Recommendation) GetUserMovieRatings(ratings []Rating) map[int64]map[int64]float64 {
	userFilmRatings := make(map[int64]map[int64]float64, 0)
	for _, rating := range ratings {
		if _, ok := userFilmRatings[rating.UserID]; !ok {
			userFilmRatings[rating.UserID] = make(map[int64]float64)
		}

		userFilmRatings[rating.UserID][rating.MovieID] = rating.Value
	}

	return userFilmRatings
}

// determine average rating
func (r *Recommendation) MeanRating(ratings []float64) float64 {
	i := len(ratings)
	if i == 0 {
		return 0
	}

	var sum float64 = 0

	for _, r := range ratings {
		sum += r
	}

	return sum / float64(i)
}

// returns: userID => [movieID => rating]
func (r *Recommendation) normalizeUserMovieRatings(userMovieRatings map[int64]map[int64]float64) {
	var ratings []float64
	var average float64
	var newRatings map[int64]float64
	for userID, movieRating := range userMovieRatings {
		if len(movieRating) == 0 {
			continue
		}
		ratings = make([]float64, 0, len(movieRating))
		for _, rating := range movieRating {
			ratings = append(ratings, rating)
		}

		average = r.MeanRating(ratings)
		newRatings = make(map[int64]float64, len(ratings))
		for movieID, rating := range movieRating {
			newRatings[movieID] = rating - average
		}

		userMovieRatings[userID] = newRatings
	}
}

// filmRatings: userID => filmRating
// userAverages: userID => userAverageRating
func (r *Recommendation) filmAverageRating(filmRatings, userAverages map[int64]float64) float64 {
	if len(filmRatings) == 0 {
		return float64(0)
	}

	ratings := make([]float64, len(filmRatings), len(filmRatings))
	i := 0
	for _, value := range filmRatings {
		ratings[i] = value
		i++
	}

	filmAverageRating := r.MeanRating(ratings)

	nominator := float64(0)
	for userID, filmRating := range filmRatings {
		nominator += filmRating - userAverages[userID]
	}

	return nominator/float64(len(filmRatings)) + filmAverageRating
}

// args: ratings1: [movieID | userID] => ratingValue
func (r *Recommendation) cosineSimilarity(ratings1, ratings2 map[int64]float64) float64 {
	var denom1, denom2 float64

	denom1 = 0
	denom2 = 0

	for _, rating := range ratings1 {
		denom1 += rating * rating
	}

	for _, rating := range ratings2 {
		denom2 += rating * rating
	}

	denominator := math.Sqrt(denom1) * math.Sqrt(denom2)
	if math.IsNaN(denominator) {
		log.Fatalf("%v, %v", ratings1, ratings2 )
	}
	if denominator == 0 {
		return 0
	}

	var nominator float64 = 0

	for id1, rating1 := range ratings1 {
		if rating2, ok := ratings2[id1]; ok {
			nominator += rating1 * rating2
		}
	}

	return nominator / denominator
}

func debug(r map[int64]map[int64]float64, onlyID bool) {
	if onlyID {
		fmt.Printf("Number of entries: %v\n", len(r))
		for id, _ := range r {
			fmt.Printf("%v, ", id)
		}
		fmt.Println()
	} else {
		for id, rates := range r {
			fmt.Printf("%v, %v\n", id, rates)
		}
	}
}