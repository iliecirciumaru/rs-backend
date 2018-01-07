package main

import (
	"log"
	"github.com/iliecirciumaru/rs-backend/db"
	"github.com/iliecirciumaru/rs-backend/repo"
	"github.com/iliecirciumaru/rs-backend/model"
)

func main() {
	//dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rsbig")
	if err != nil {
		log.Fatal(err)
	}

	ratingRepo := repo.NewRatingRepo(dbsess)

	mostRatedMovies, err := ratingRepo.GetMaxNumberRatedMovies(1)
	if err != nil {
		log.Fatal(err)
	}

	rec := model.Recommendation{UuNeighbours: 3}
	ca := model.ClusteringUtility{Rec: rec}

	ratings, err := ratingRepo.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	ca.Cluster(ratings, mostRatedMovies[0].Key)

}
