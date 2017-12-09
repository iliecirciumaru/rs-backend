package model

type Movie struct {
	ID          int64  `db:"id,omitempty"`
	Title       string `db:"title"`
	Information string `db:"information"`
	PosterURL   *string `db:"poster_image_url"`
}

type UserMovieView struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Information string  `json:"info"`
	UserRating  float64 `json:"user_rating"`
	PosterURL   string  `json:"poster_image_url"`
}

type MovieView struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	Information   string  `json:"info"`
	AverageRating float64 `json:"average_rating"`
	PosterURL     string  `json:"poster_image_url"`
}

type MoviePrediction struct {
	MovieID int64
	PredictedScore float64
}

type ByScoreDesc []MoviePrediction

func (a ByScoreDesc) Len() int           { return len(a) }
func (a ByScoreDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScoreDesc) Less(i, j int) bool { return a[i].PredictedScore > a[j].PredictedScore }