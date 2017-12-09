package main

import (
	"encoding/json"
	"fmt"
	"github.com/iliecirciumaru/rs-backend/db"
	"github.com/iliecirciumaru/rs-backend/model"
	"io/ioutil"
	"log"
	"math"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"github.com/iliecirciumaru/rs-backend/evaluation/structs"
)

func main() {
	dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	if err != nil {
		log.Fatal(err)
	}

	neighbours := []uint{5, 10, 15, 20, 25, 30}
	uuResults := make([]structs.EvaluationUUResult, len(neighbours))

	for i, n := range neighbours {
		uuResults[i] = EvaluateUURecommender(dbsess, n)
	}

	rawResult, _ := json.MarshalIndent(uuResults, "", "    ")
	saveResults(rawResult, "uuCLF")
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

func EvaluateUURecommender(dbsess sqlbuilder.Database, neighbours uint) structs.EvaluationUUResult {
	uuResult := structs.EvaluationUUResult{Neighbours: neighbours}

	recommender := model.Recommendation{neighbours}

	// accumulator for goabal rmse
	globalRMSE := float64(0)
	globalCount := 0

	//userRMSE := float64(0)
	//userCount := 0

	userPredictedCount := 0

	testratings := getRatings(dbsess, "testratings")
	testUserMovieRating := recommender.GetUserMovieRatings(testratings)

	normalRatings := getRatings(dbsess, "ratings")
	userMovieRating := recommender.GetUserMovieRatings(normalRatings)

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
		scores := recommender.PredictUserScores(userID, testratings)
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
	uuResult.TotalTime = float64(end - start) / 1000000000
	uuResult.TimePerUser = float64(uuResult.TotalTime) / float64(userPredictedCount)
	return uuResult
}
