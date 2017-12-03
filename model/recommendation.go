package model

type Recommendation struct {

}


func (r *Recommendation) ComputeFilmAverages(ratings []Rating) map[int64]float64 {
	userRatings := make(map[int64][]float64, 0)
	filmRatings := make(map[int64]map[int64]float64, 0)

	for _, rating := range ratings {
		if _, ok := userRatings[rating.UserID]; ok {
			userRatings[rating.UserID] = append(userRatings[rating.UserID], rating.Value)
		} else {
			userRatings[rating.UserID] = make([]float64, 1)
			userRatings[rating.UserID][0] = rating.Value
		}

		if _, ok := filmRatings[rating.MovieID]; !ok {
			filmRatings[rating.MovieID] = make(map[int64]float64)
		}

		filmRatings[rating.MovieID][rating.UserID] = rating.Value
	}

	userAverages := make(map[int64]float64)
	for userID, userRats := range userRatings {
		userAverages[userID] = r.MeanRating(userRats)
	}


	filmAverages := make(map[int64]float64)
	for filmID, userFilmRatings := range filmRatings {
		if len(userFilmRatings) < 10 {
			filmAverages[filmID] = 0
			continue
		}

		filmAverages[filmID] = r.filmAverageRating(userFilmRatings, userAverages)
	}


	return filmAverages
}

// determine average rating
func (r* Recommendation) MeanRating(ratings []float64) float64 {
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

// filmRatings: userID => filmRating
// userAverages: userID => userAverageRating
func (r* Recommendation) filmAverageRating(filmRatings, userAverages map[int64]float64) float64 {
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

	return nominator / float64(len(filmRatings)) + filmAverageRating
}