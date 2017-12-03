package repo

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"github.com/iliecirciumaru/rs-backend/model"
	"fmt"
)

func NewRatingRepo(db sqlbuilder.Database) RatingRepo {
	return RatingRepo{
		db: db,
	}
}

type RatingRepo struct {
	db sqlbuilder.Database
}

func (r *RatingRepo) AddRating(rating model.Rating) error {
	_, err := r.db.Collection("ratings").Insert(&rating)

	return err
}

func (r *RatingRepo) GetRatingByMovieUserID(userID, movieID int64) model.Rating {
	var rating model.Rating
	query := r.db.SelectFrom("ratings").Where("iduser = ? and idmovie = ?", userID, movieID)
	err := query.One(&rating)
	if err != nil {
		fmt.Printf("RatingRepo, err: %s\n", err)
	}
	fmt.Println(rating)

	return rating
}

func (r *RatingRepo) GetAll() ([]model.Rating, error) {
	var ratings []model.Rating
	err := r.db.SelectFrom("ratings").All(&ratings)

	return ratings, err
}