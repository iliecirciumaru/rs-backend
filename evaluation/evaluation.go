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
)

var recommender model.Recommendation = model.Recommendation{UuNeighbours:2}

func main() {
	//dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rsbig")
	if err != nil {
		log.Fatal(err)
	}

	neighbours := []uint{2, 4, 5, 6, 8, 10, 12, 14, 15, 16, 20, 25, 30}
	//uuResults := make([]structs.EvaluationResult, len(neighbours))


	//// test UUCLF recommender
	//for i, n := range neighbours {
	//	uuResults[i] = EvaluateRecommender(dbsess, n, true)
	//}
	//
	//rawResult, _ := json.MarshalIndent(uuResults, "", "    ")
	//saveResults(rawResult, "uuCLF")


	// test IICLF recommender
	//neighbours = []uint{2, 4, 5, 6, 8, 10, 12, 14, 15, 16, 20, 25, 30}
	neighbours = []uint{15}
	iiResults := make([]structs.EvaluationResult, len(neighbours))

	for i, n := range neighbours {
		iiResults[i] = EvaluateRecommender(dbsess, n, false)
	}

	rawResult, _ := json.MarshalIndent(iiResults, "", "    ")
	saveResults(rawResult, "iiCLF")

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

	//testratings := getRatings(dbsess, "testratings")
	percent := int(float32(len(normalRatings))*0.8)
	testratings := normalRatings[0: percent]
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
