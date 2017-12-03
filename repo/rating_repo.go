package repo

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"github.com/iliecirciumaru/rs-backend/model"
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