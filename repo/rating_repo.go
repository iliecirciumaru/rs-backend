package repo

import (
	"fmt"
	"github.com/iliecirciumaru/rs-backend/model"
	"upper.io/db.v3/lib/sqlbuilder"
	"strconv"
	"github.com/iliecirciumaru/rs-backend/structs"
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

func (r *RatingRepo) GetAVGMovieRating(movieIDs []int64) ([]structs.KeyValue, error){
	var ids string = ""

	for i, id := range movieIDs {
		if i != 0 {
			ids += ","
		}
		ids += strconv.FormatInt(id, 10)
	}

	var res []structs.KeyValue

	qs := fmt.Sprintf("SELECT idmovie as 'key', AVG(rating) as 'value' FROM ratings WHERE idmovie IN (%s) GROUP BY idmovie", ids)

	rows, err := r.db.Query(qs)
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&res)

	if err != nil {
		fmt.Println(err)
	}

	return res, err
}

// return movies which are rated the most
func (r *RatingRepo) GetMaxNumberRatedMovies(number int) ([]structs.KeyValue, error){
	var res []structs.KeyValue

	qs := fmt.Sprintf("SELECT idmovie as 'key', COUNT(idmovie) as 'value' " +
		"FROM ratings GROUP BY idmovie ORDER BY value DESC LIMIT %d", number)

	rows, err := r.db.Query(qs)
	if err != nil {
		return nil, err
	}

	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&res)

	if err != nil {
		fmt.Println(err)
	}

	return res, err
}

