package repo

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"github.com/iliecirciumaru/rs-backend/model"
	"fmt"
	"strconv"
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

func (r *MovieRepo) GetMovieByID(movieID int64) (model.Movie, error) {
	var movie model.Movie
	query := r.db.SelectFrom("movies").Where("id = ?", movieID)
	err := query.One(&movie)

	return movie, err
}

func (r *MovieRepo) GetMovieByIDs(movieIDs []int64) ([]model.Movie, error) {
	var ids string = ""

	for i, id := range movieIDs {
		if i != 0 {
			ids += ","
		}
		ids += strconv.FormatInt(id, 10)
	}

	fmt.Println(ids)


	var movies []model.Movie

	rows, err := r.db.Query(fmt.Sprintf("SELECT * FROM movies WHERE id IN (%s)", ids))
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&movies)
	//query := r.db.SelectFrom("movies").Where("id IN (?)", ids)

	//err := query.All(&movies)
	//query.
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(movies)

	return movies, err
}