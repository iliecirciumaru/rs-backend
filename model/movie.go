package model

type Movie struct {
	ID int64 `db:"id"`
	Title string `db:"title"`
	Information string `db:"information"`
}

type MovieView struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	Information string `json:"info"`
	UserRating float64 `json:"user_rating"`
}