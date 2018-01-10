package main

import (
	"log"
	"github.com/iliecirciumaru/rs-backend/db"
	"github.com/iliecirciumaru/rs-backend/repo"
	"github.com/iliecirciumaru/rs-backend/model"
	"fmt"
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
	ca := model.ClusteringUtility{
		Rec: rec,
		MinCentroidRates: 200,
		ClusterNum: 6,
	}

	ratings, err := ratingRepo.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	clusters := ca.Cluster(ratings[0:int(float32(len(ratings))*0.8)], mostRatedMovies[0].Key)

	var filteredRatings []model.Rating
	for centroid, c := range clusters {
		filteredRatings = ca.ExtractRatings(ratings, c)
		fmt.Printf("Update Similarities, cluster %v, ratings %v\n", centroid, len(filteredRatings))
		//recommender.UpdateSimilarities(filteredRatings)
	}

}
