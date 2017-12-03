package repo

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"github.com/iliecirciumaru/rs-backend/model"
)

func NewMovieRepo(db sqlbuilder.Database) MovieRepo {
	return MovieRepo{
		db: db,
	}
}

type MovieRepo struct {
	db sqlbuilder.Database
}

//func (r *RatingRepo) AddRating(rating model.Rating) error {
//	_, err := r.db.Collection("ratings").Insert(&rating)
//
//	return err
//}

func (r *MovieRepo) GetMovie(movieID int64) (model.Movie, error) {
	var movie model.Movie
	query := r.db.SelectFrom("movies").Where("id = ? ", movieID)
	err := query.One(&movie)

	return movie, err
}