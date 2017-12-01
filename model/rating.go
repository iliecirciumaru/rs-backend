package model

type Rating struct {
	UserID int `json:"user_id"`
	MovieID int `json:"user_id"`
	Timestamp int `json:"timestamp"`
	Value float64 `json:"rating"`
}