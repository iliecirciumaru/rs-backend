package model

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sort"
	"sync"
	"time"
)

type Recommendation struct {
	UuNeighbours     uint
	MovieSimilarties map[int64][]Similarity
	userMovieRatings map[int64]map[int64]float32
	movieUserRatings map[int64]map[int64]float32
	movieAverages    map[int64]float32
}

// returns: movieID => averageScore (Normalized)
func (r *Recommendation) ComputeFilmAverageScores(ratings []Rating) map[int64]float32 {
	userRatings := r.getUserRatingVectors(ratings)
	filmRatings := r.getMovieUserRatings(ratings)

	userAverages := make(map[int64]float32)
	for userID, userRats := range userRatings {
		userAverages[userID] = r.MeanRating(userRats)
	}

	filmScores := make(map[int64]float32)
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
func (r *Recommendation) PredictUserScoreUUCLF(userID int64, ratings []Rating) []MoviePrediction {
	userMovieRatings := r.GetUserMovieRatings(ratings)
	userAverageScore := float32(0)

	// user hasn't rated any movies, we don't predict anything
	if _, ok := userMovieRatings[userID]; !ok {
		return nil
	}

	for _, rating := range userMovieRatings[userID] {
		userAverageScore += rating
	}

	userAverageScore /= float32(len(userMovieRatings[userID]))

	r.normalizeUserMovieOrMovieUserRatings(userMovieRatings)

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

	var score float32
	for movieID, uRatings := range movieUserRatings {
		// if user has rated this movie, we continue
		if _, ok := uRatings[userID]; ok {
			continue
		}

		// take into consideration only top n neighbours
		neighbours := uint(0)
		nominator := float32(0)
		denominator := float32(0)
		for _, similarity := range cosineSimilarities {

			if neighbours == r.UuNeighbours {
				score = nominator/denominator + userAverageScore
				if math.IsNaN(float64(score)) {
					fmt.Println("NaN score for movie %v, user %v\n", movieID, userID)
				}
				if score > 8 {
					fmt.Printf("Movie %v, User %v\n, Score %v, denom %v, nomin %v\n", movieID, userID, score, denominator, nominator)
				}

				filmPredictions = append(filmPredictions, MoviePrediction{MovieID: movieID, PredictedScore: score})
				break
			}

			if _, ok := uRatings[similarity.ID]; ok {
				neighbours++
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

// returns: movieID => predictedScore
// are returned only movies, which user hasn't rated
func (r *Recommendation) PredictUserScoreIICLF(userID int64, ratings []Rating) []MoviePrediction {
	if r.userMovieRatings == nil {
		r.userMovieRatings = r.GetUserMovieRatings(ratings)
	}

	// user hasn't rated any movies, we don't predict anything
	if _, ok := r.userMovieRatings[userID]; !ok {
		return nil
	}

	uRatings := r.userMovieRatings[userID]

	if r.movieUserRatings == nil {
		r.movieUserRatings = r.getMovieUserRatings(ratings)
		r.movieAverages = make(map[int64]float32, len(r.movieUserRatings))
		for movieID, userRatings := range r.movieUserRatings {
			rates := make([]float32, len(userRatings))
			i := 0
			for _, r := range userRatings {
				rates[i] = r
				i++
			}

			r.movieAverages[movieID] = r.MeanRating(rates)
		}

		r.normalizeUserMovieOrMovieUserRatings(r.movieUserRatings)
	}




	filmPredictions := make([]MoviePrediction, 0, 20)

	if r.MovieSimilarties == nil {
		fmt.Println("Calculate Similarity from IICLF")
		r.calculateMovieSimilarites(r.movieUserRatings)
	}

	// make predictions
	for movieID, _ := range r.movieUserRatings {
		// if user has rated this movie, we continue
		if _, ok := uRatings[movieID]; ok {
			continue
		}

		denominator := float32(0)
		nominator := float32(0)
		neighbours := uint(0)
		score := float32(0)

		for _, similarity := range r.MovieSimilarties[movieID] {

			if neighbours == r.UuNeighbours {
				score = nominator/denominator + r.movieAverages[movieID]
				if math.IsNaN(float64(score)) {
					fmt.Println("NaN score for movie %v, user %v\n", movieID, userID)
				}
				if score > 8 {
					fmt.Printf("Movie %v, User %v\n, Score %v, denom %v, nomin %v\n", movieID, userID, score, denominator, nominator)
				}
				//fmt.Printf("Score for film %v, equals %v\n", movieID, score)

				filmPredictions = append(filmPredictions, MoviePrediction{MovieID: movieID, PredictedScore: score})
				break
			}

			if movieRating, ok := uRatings[similarity.ID]; ok {
				neighbours++
				// if positive similarities are not enough
				// we don't recommend this film at all
				if similarity.Value <= 0 {
					break
				}

				denominator += similarity.Value
				nominator += similarity.Value * (movieRating - r.movieAverages[similarity.ID])
			}
		}
	}

	sort.Sort(ByScoreDesc(filmPredictions))

	return filmPredictions
}

func (r *Recommendation) calculateMovieSimilarites(movieUserRatings map[int64]map[int64]float32) {
	fmt.Println("Start calculating movie similarities")
	start := time.Now().UnixNano()
	var wg sync.WaitGroup
	mutex := sync.Mutex{}

	if r.MovieSimilarties == nil {
		r.MovieSimilarties = make(map[int64][]Similarity)
	}

	//for movieID, uRatings := range movieUserRatings {
	//	wg.Add(1)
	//	go r.CalculateMovieSimilarity(movieID, uRatings, &wg, &mutex, movieUserRatings)
	//}

	jobs := make(chan int64, 100)
	//results := make(chan int64, 100)

	// launch worker per cpu core
	for i := 1; i <= runtime.NumCPU(); i++ {
		go r.CalculateMovieSimilarity(jobs, &wg, &mutex, movieUserRatings)
	}

	j := int64(0)

	for movieID, _ := range movieUserRatings {
		wg.Add(1)
		jobs <- movieID
		j++
		if j%1000 == 0 {
			fmt.Println(j/1000, "k")
		}
	}
	close(jobs)

	wg.Wait()
	end := time.Now().UnixNano()
	fmt.Printf("Movie similarities are calculated and cached, time: %.2fs\n", float32(end-start)/1000000000)
}

func (r *Recommendation) UpdateSimilarities(ratings []Rating) {
	r.calculateMovieSimilarites(r.getMovieUserRatings(ratings))
}

func (r *Recommendation) GetMostSimilarMovies(movieID int64, ratings []Rating) []Similarity {
	movieUserRatings := r.getMovieUserRatings(ratings)
	if _, ok := movieUserRatings[movieID]; !ok {
		return nil
	}

	r.normalizeUserMovieOrMovieUserRatings(movieUserRatings)
	cosineSimilarities := make([]Similarity, 0, len(movieUserRatings)-1)

	userRatings := movieUserRatings[movieID]
	// calculate cosine similarites
	for movieID2, userRatings2 := range movieUserRatings {
		if movieID == movieID2 {
			continue
		}

		cosineSimilarities = append(cosineSimilarities, Similarity{
			ID:    movieID2,
			Value: r.cosineSimilarity(userRatings, userRatings2),
		})

	}

	sort.Sort(BySimilarityDesc(cosineSimilarities))

	return cosineSimilarities
}

// return: userID => [movieRatingi, movieRating]
func (r *Recommendation) getUserRatingVectors(ratings []Rating) map[int64][]float32 {
	userRatings := make(map[int64][]float32, 0)
	for _, rating := range ratings {
		if _, ok := userRatings[rating.UserID]; ok {
			userRatings[rating.UserID] = append(userRatings[rating.UserID], rating.Value)
		} else {
			userRatings[rating.UserID] = make([]float32, 1)
			userRatings[rating.UserID][0] = rating.Value
		}
	}

	return userRatings
}

// returns: movieID => [userID1 => rating, userID2 => rating]
func (r *Recommendation) getMovieUserRatings(ratings []Rating) map[int64]map[int64]float32 {
	filmRatings := make(map[int64]map[int64]float32, 0)
	for _, rating := range ratings {
		if _, ok := filmRatings[rating.MovieID]; !ok {
			filmRatings[rating.MovieID] = make(map[int64]float32)
		}

		filmRatings[rating.MovieID][rating.UserID] = rating.Value
	}

	return filmRatings
}

// returns: userID => [movieID1 => rating, movieID2 => rating]
func (r *Recommendation) GetUserMovieRatings(ratings []Rating) map[int64]map[int64]float32 {
	userFilmRatings := make(map[int64]map[int64]float32, 0)
	for _, rating := range ratings {
		if _, ok := userFilmRatings[rating.UserID]; !ok {
			userFilmRatings[rating.UserID] = make(map[int64]float32)
		}

		userFilmRatings[rating.UserID][rating.MovieID] = rating.Value
	}

	return userFilmRatings
}

// determine average rating
func (r *Recommendation) MeanRating(ratings []float32) float32 {
	i := len(ratings)
	if i == 0 {
		return 0
	}

	var sum float32 = 0

	for _, r := range ratings {
		sum += r
	}

	return sum / float32(i)
}

// returns: userID => [movieID => rating] || movieID => [userID => rating]
func (r *Recommendation) normalizeUserMovieOrMovieUserRatings(IDtoIDratings map[int64]map[int64]float32) {
	var ratings []float32
	var average float32
	var newRatings map[int64]float32
	for userID, movieRating := range IDtoIDratings {
		if len(movieRating) == 0 {
			continue
		}
		ratings = make([]float32, 0, len(movieRating))
		for _, rating := range movieRating {
			ratings = append(ratings, rating)
		}

		average = r.MeanRating(ratings)
		newRatings = make(map[int64]float32, len(ratings))
		for movieID, rating := range movieRating {
			newRatings[movieID] = rating - average
		}

		IDtoIDratings[userID] = newRatings
	}
}

// filmRatings: userID => filmRating
// userAverages: userID => userAverageRating
func (r *Recommendation) filmAverageRating(filmRatings, userAverages map[int64]float32) float32 {
	if len(filmRatings) == 0 {
		return float32(0)
	}

	ratings := make([]float32, len(filmRatings), len(filmRatings))
	i := 0
	for _, value := range filmRatings {
		ratings[i] = value
		i++
	}

	filmAverageRating := r.MeanRating(ratings)

	nominator := float32(0)
	for userID, filmRating := range filmRatings {
		nominator += filmRating - userAverages[userID]
	}

	return nominator/float32(len(filmRatings)) + filmAverageRating
}

// args: ratings1: [movieID | userID] => ratingValue
func (r *Recommendation) cosineSimilarity(ratings1, ratings2 map[int64]float32) float32 {
	var denom1, denom2 float32

	denom1 = 0
	denom2 = 0

	for _, rating := range ratings1 {
		denom1 += rating * rating
	}

	for _, rating := range ratings2 {
		denom2 += rating * rating
	}

	denominator := math.Sqrt(float64(denom1)) * math.Sqrt(float64(denom2))
	if math.IsNaN(denominator) {
		log.Fatalf("%v, %v", ratings1, ratings2)
	}
	if denominator == 0 {
		return 0
	}

	var nominator float32 = 0

	for id1, rating1 := range ratings1 {
		if rating2, ok := ratings2[id1]; ok {
			nominator += rating1 * rating2
		}
	}

	return float32(float64(nominator) / denominator)
}

func debug(r map[int64]map[int64]float32, onlyID bool) {
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
