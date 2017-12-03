package model

type Rating struct {
	UserID int `json:"user_id" db:"iduser"`
	MovieID int `json:"user_id" db:"idmovie"`
	Timestamp int64 `json:"timestamp" db:"timestamp"`
	Value float64 `json:"rating" db:"rating"`
}


type RatingAddRequest struct {
	MovieID int `json:"movie_id"`
	Value float64 `json:"rating"`
}