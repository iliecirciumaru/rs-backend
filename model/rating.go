package model

type Rating struct {
	UserID    int64   `json:"user_id" db:"iduser"`
	MovieID   int64   `json:"user_id" db:"idmovie"`
	Timestamp int64   `json:"timestamp" db:"timestamp"`
	Value     float64 `json:"rating" db:"rating"`
}

type RatingAddRequest struct {
	MovieID int64   `json:"movie_id"`
	Value   float64 `json:"rating"`
}

func RoundRating(rating float64) float64 {
	return float64(int64(rating*20+0.5)) / 20
}