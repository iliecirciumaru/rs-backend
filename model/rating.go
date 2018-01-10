package model

type Rating struct {
	UserID    int64   `json:"user_id" db:"iduser"`
	MovieID   int64   `json:"user_id" db:"idmovie"`
	Timestamp int64   `json:"timestamp" db:"timestamp"`
	Value     float32 `json:"rating" db:"rating"`
}

type RatingAddRequest struct {
	MovieID int64   `json:"movie_id"`
	Value   float32 `json:"rating"`
}

func RoundRating(rating float32) float32 {
	return float32(int64(rating*20+0.5)) / 20
}