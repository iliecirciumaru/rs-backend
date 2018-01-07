package main

import (
	"encoding/json"
	"fmt"
	"github.com/iliecirciumaru/rs-backend/db"
	"github.com/iliecirciumaru/rs-backend/evaluation/structs"
	"github.com/iliecirciumaru/rs-backend/model"
	"io/ioutil"
	"log"
	"math"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"github.com/iliecirciumaru/rs-backend/repo"
)

var recommender model.Recommendation = model.Recommendation{UuNeighbours:2}
var clusterUtility model.ClusteringUtility = model.ClusteringUtility{recommender}

func main() {
	//dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rsbig")
	if err != nil {
		log.Fatal(err)
	}

	//evaluateUUCLF(dbsess)
	//evaluateIICLF(dbsess)
	evaluateIICLFClustering(dbsess)

}

func evaluateUUCLF(dbsess sqlbuilder.Database) {
	//neighbours := []uint{2, 4, 5, 6, 8, 10, 12, 14, 15, 16, 20, 25, 30}
	neighbours := []uint{5, 15,25}
	uuResults := make([]structs.EvaluationResult, len(neighbours))


	// test UUCLF recommender
	for i, n := range neighbours {
		uuResults[i] = EvaluateRecommender(dbsess, n, true)
	}

	rawResult, _ := json.MarshalIndent(uuResults, "", "    ")
	saveResults(rawResult, "uuCLF")
}


func evaluateIICLF(dbsess sqlbuilder.Database) {
	// test IICLF recommender
	neighbours := []uint{2, 4, 5, 6, 8, 10, 12, 14, 15, 16, 20, 25, 30}
	//neighbours := []uint{15}
	iiResults := make([]structs.EvaluationResult, len(neighbours))

	for i, n := range neighbours {
		iiResults[i] = EvaluateRecommender(dbsess, n, false)
	}

	rawResult, _ := json.MarshalIndent(iiResults, "", "    ")
	saveResults(rawResult, "iiCLF")
}

func evaluateIICLFClustering(dbsess sqlbuilder.Database) {
	ratingRepo := repo.NewRatingRepo(dbsess)
	mostRatedMovies, err := ratingRepo.GetMaxNumberRatedMovies(1)
	if err != nil {
		log.Fatal(err)
	}

	ratings := getRatings(dbsess, "ratings")

	clusters := clusterUtility.Cluster(ratings[0: getTopRating(len(ratings))], mostRatedMovies[0].Key)
	var filteredRatings []model.Rating
	for _, c := range clusters {
		filteredRatings = clusterUtility.ExtractRatings(ratings, c)
		recommender.UpdateSimilarities(filteredRatings)
	}

	//neighbours := []uint{2, 4, 5, 6, 8, 10, 12, 14, 15, 16, 20, 25, 30}
	neighbours := []uint{5, 10, 15, 20}
	iiResults := make([]structs.EvaluationResult, len(neighbours))

	for i, n := range neighbours {
		iiResults[i] = EvaluateRecommender(dbsess, n, false)
	}

	rawResult, _ := json.MarshalIndent(iiResults, "", "    ")
	saveResults(rawResult, "iiCLF_cluster_big")
}

func saveResults(data []byte, filename string) {
	ioutil.WriteFile(filename+".json", data, 0644)
}

func getRatings(dbsses sqlbuilder.Database, tableName string) []model.Rating {
	var ratings []model.Rating
	err := dbsses.SelectFrom(tableName).All(&ratings)
	if err != nil {
		log.Fatal(err)
	}

	return ratings
}

func EvaluateRecommender(dbsess sqlbuilder.Database, neighbours uint, uuCLF bool) structs.EvaluationResult {
	uuResult := structs.EvaluationResult{Neighbours: neighbours}

	recommender.UuNeighbours = neighbours

	// accumulator for goabal rmse
	globalRMSE := float64(0)
	globalCount := 0

	//userRMSE := float64(0)
	//userCount := 0

	userPredictedCount := 0

	normalRatings := getRatings(dbsess, "ratings")
	userMovieRating := recommender.GetUserMovieRatings(normalRatings)

	testratings := getRatings(dbsess, "testratings")
	//testratings := normalRatings[0: getTopRating(len(normalRatings))]
	testUserMovieRating := recommender.GetUserMovieRatings(testratings)

	start := time.Now().UnixNano()

	for userID, testMovieRatings := range testUserMovieRating {
		if len(testMovieRatings) == 0 {
			continue
		}
		//if userID != 1 {
		//	continue
		//}
		userPredictedCount++
		if userPredictedCount > 1000 {
			break
		}

		//userRMSE = 0
		//userCount = 0

		//fmt.Printf("Start prediction for user %v\n", userID)
		var scores []model.MoviePrediction
		if uuCLF {
			scores = recommender.PredictUserScoreUUCLF(userID, testratings)
		} else {
			scores = recommender.PredictUserScoreIICLF(userID, testratings)
		}

		//fmt.Printf("Predicted scores for %d movies\n", len(scores))
		for _, prediction := range scores {

			if realRating, ok := userMovieRating[userID][prediction.MovieID]; ok {
				globalRMSE += (prediction.PredictedScore - realRating) * (prediction.PredictedScore - realRating)
				globalCount++

				//userRMSE += (prediction.PredictedScore - realRating) * (prediction.PredictedScore - realRating)
				//userCount++
			}
		}

		//if userCount > 0 {
		//	fmt.Printf("RMSE for user %v equals: %5.2f\n", userID, math.Sqrt(userRMSE/float64(userCount)))
		//}
	}

	end := time.Now().UnixNano()

	if globalCount == 0 {
		log.Fatal("Count is zero, smth went wrong")
	}
	globalRMSE = math.Sqrt(globalRMSE / float64(globalCount))

	fmt.Printf("Neighbours: %v, RMSE equals to: %5.2f, predicted scores %d\n", neighbours, globalRMSE, globalCount)

	uuResult.GlobalRMSE = globalRMSE
	uuResult.PredictionsUsed = globalCount
	uuResult.TotalTime = float64(end-start) / 1000000000
	uuResult.TimePerUser = float64(uuResult.TotalTime) / float64(userPredictedCount)
	return uuResult
}


func getTopRating(allRatingLength int) int {
	return int(float32(allRatingLength)*0.8)
}

func getBottomRating(allRatingLength int) int {
	return int(float32(allRatingLength)*0.2)
}
