package repo

import (
	"fmt"
	"github.com/iliecirciumaru/rs-backend/model"
	"strconv"
	"upper.io/db.v3/lib/sqlbuilder"
)

func NewMovieRepo(db sqlbuilder.Database) MovieRepo {
	return MovieRepo{
		db: db,
	}
}

type MovieRepo struct {
	db sqlbuilder.Database
}

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

	var movies []model.Movie

	rows, err := r.db.Query(fmt.Sprintf("SELECT * FROM movies WHERE id IN (%s)", ids))
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&movies)

	if err != nil {
		fmt.Println(err)
	}

	return movies, err
}

func (r *MovieRepo) GetLatestMovies(number int) ([]model.Movie, error) {
	var movies []model.Movie
	// TODO change to desc, when latest movies have poster urls
	rows, err := r.db.Query(fmt.Sprintf("SELECT * FROM movies ORDER BY id DESC LIMIT %d", number))

	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&movies)

	if err != nil {
		fmt.Println(err)
	}

	return movies, err
}

func (r *MovieRepo) GetMoviesByPrefix(prefix string) ([]model.Movie, error) {
	var movies []model.Movie
	sqlStr := fmt.Sprintf("SELECT * FROM movies WHERE title LIKE '%s%%' ORDER BY id ASC LIMIT 5", prefix)
	fmt.Println(sqlStr)

	rows, err := r.db.Query(sqlStr)

	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&movies)

	if err != nil {
		fmt.Println(err)
	}

	return movies, err
}

func (r *MovieRepo) GetAll() ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.SelectFrom("movies").All(&movies)



	return movies, err
}