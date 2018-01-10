package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/iliecirciumaru/rs-backend/db"
	"github.com/iliecirciumaru/rs-backend/evaluation/structs"
	"github.com/iliecirciumaru/rs-backend/model"
	"github.com/iliecirciumaru/rs-backend/repo"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"math/rand"
)

var recommender model.Recommendation = model.Recommendation{UuNeighbours: 2}
var clusterUtility model.ClusteringUtility = model.ClusteringUtility{
	Rec:              recommender,
	MinCentroidRates: 250,
	ClusterNum: 4,
}

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
	neighbours := []uint{5, 15, 25}
	uuResults := make([]structs.EvaluationResult, len(neighbours))

	// test UUCLF recommender
	for i, n := range neighbours {
		uuResults[i] = EvaluateRecommender(dbsess, n, true, nil, nil)
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
		iiResults[i] = EvaluateRecommender(dbsess, n, false, nil, nil)
	}

	rawResult, _ := json.MarshalIndent(iiResults, "", "    ")
	saveResults(rawResult, "iiCLF")
}

func saveMovieSimilarities(filename string) error {
	file := filename + "_movie_similarities"
	fmt.Printf("Start saving movie similarities into file %s\n", file)
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)

	err := e.Encode(recommender.MovieSimilarties)
	if err != nil {
		fmt.Printf("Error during encoding: %s\n", err.Error())
		return err
	}

	err = ioutil.WriteFile(file, b.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error during writing to file: %s\n", err.Error())
		return err
	}

	fmt.Println("Saving of similarities was succesfully done")
	return nil
}

func extractMovieSimilaties(filename string) error {
	file := filename + "_movie_similarities"
	fmt.Printf("Start extracting movie similarities into file %s\n", file)
	if _, err := os.Stat(file); err != nil {
		fmt.Printf("No file %s, error: %s\n", filename, err.Error())
		return err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Canr read file %s, error: %s\n", filename, err.Error())
		return err
	}

	b := bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	err = d.Decode(&recommender.MovieSimilarties)
	if err != nil {
		fmt.Printf("Can't decode similarities, err %s\n", err.Error())
		return err
	}

	fmt.Println("Extracting of similarities was succesfully done")
	return nil
}

func evaluateIICLFClustering(dbsess sqlbuilder.Database) {
	ratingRepo := repo.NewRatingRepo(dbsess)

	ratings := getRatings(dbsess, "ratings")

	//if extractMovieSimilaties("cluter_big_sims") != nil {
	mostRatedMovies, err := ratingRepo.GetMaxNumberRatedMovies(1)
	if err != nil {
		log.Fatal(err)
	}

	//testRatings := ratings[0:getTopRating(len(ratings))]
	//testRatings := getRatings(dbsess, "testratings")
	testRatings := getTestRatings(ratings)
	fmt.Println("Test rating length: ", len(testRatings))
	clusters := clusterUtility.Cluster(testRatings, mostRatedMovies[0].Key)

	var filteredRatings []model.Rating
	for centroid, c := range clusters {
		filteredRatings = clusterUtility.ExtractRatings(testRatings, c)
		fmt.Printf("Update Similarities, cluster %v, ratings %v\n", centroid, len(filteredRatings))
		recommender.UpdateSimilarities(filteredRatings)
	}

	//saveMovieSimilarities("cluter_big_sims")
	//}

	//neighbours := []uint{2, 4, 5, 6, 8, 10, 12, 14, 15, 16, 20, 25, 30}
	neighbours := []uint{3, 5, 15}
	iiResults := make([]structs.EvaluationResult, len(neighbours))

	for i, n := range neighbours {
		iiResults[i] = EvaluateRecommender(dbsess, n, false, ratings, testRatings)
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

func EvaluateRecommender(dbsess sqlbuilder.Database, neighbours uint, uuCLF bool, ratings, testratings []model.Rating) structs.EvaluationResult {
	fmt.Printf("Start recommender evaluation, neighbours: %v, uuCLF: %v\n", neighbours, uuCLF)
	uuResult := structs.EvaluationResult{Neighbours: neighbours}
	usersRepo := repo.NewUserRepo(dbsess)
	users, err := usersRepo.GetAll()
	if err != nil {
		log.Fatalf("Can't retrieve users")
	}

	recommender.UuNeighbours = neighbours

	// accumulator for goabal rmse
	globalRMSE := float64(0)
	globalCount := 0

	//userRMSE := float64(0)
	//userCount := 0

	userPredictedCount := 0

	if ratings == nil {
		ratings = getRatings(dbsess, "ratings")
	}

	userMovieRating := recommender.GetUserMovieRatings(ratings)

	if testratings == nil {
		testratings = ratings[0:getTopRating(len(ratings))]
	}

	//testUserMovieRating := recommender.GetUserMovieRatings(testratings)

	start := time.Now().UnixNano()

	for _, u := range users {
		u = users[int(rand.Uint32())%len(users)]
		var scores []model.MoviePrediction
		if uuCLF {
			scores = recommender.PredictUserScoreUUCLF(u.ID, testratings)
		} else {
			scores = recommender.PredictUserScoreIICLF(u.ID, testratings)
		}

		if scores == nil {
			continue
		}

		userPredictedCount++
		if userPredictedCount > 1500 {
			break
		}

		if userPredictedCount%100 == 0 {
			fmt.Println(userPredictedCount/100, "h")
		}

		fmt.Printf("Predicted for %d user, %d scores\n", u.ID, len(scores))
		for _, prediction := range scores {

			if realRating, ok := userMovieRating[u.ID][prediction.MovieID]; ok {
				globalRMSE += float64((prediction.PredictedScore - realRating) * (prediction.PredictedScore - realRating))
				globalCount++

				//userRMSE += (prediction.PredictedScore - realRating) * (prediction.PredictedScore - realRating)
				//userCount++
			}
		}
		fmt.Println("Global count: ", globalCount)

		//if userCount > 0 {
		//	fmt.Printf("RMSE for user %v equals: %5.2f\n", userID, math.Sqrt(userRMSE/float64(userCount)))
		//}
	}

	end := time.Now().UnixNano()

	if globalCount == 0 {
		fmt.Println("NO RATESSSSS, count is zero")
		//log.Fatal("Count is zero, smth went wrong")
		globalCount = 1
	}
	globalRMSE = math.Sqrt(globalRMSE / float64(globalCount))

	uuResult.GlobalRMSE = globalRMSE
	uuResult.PredictionsUsed = globalCount
	uuResult.TotalTime = float64(end-start) / 1000000000
	uuResult.TimePerUser = float64(uuResult.TotalTime) / float64(userPredictedCount)

	fmt.Printf("Neighbours: %v, RMSE equals to: %5.2f, predicted scores %d, time:%.2f\n", neighbours, globalRMSE, globalCount, uuResult.TotalTime)
	return uuResult
}

func getTopRating(allRatingLength int) int {
	return int(float64(allRatingLength) * 0.75)
}

func getBottomRating(allRatingLength int) int {
	return int(float64(allRatingLength) * 0.4)
}

func getTestRatings(ratings []model.Rating) []model.Rating {
	size := int(float64(len(ratings)) * 0.78)

	testRatings := make([]model.Rating, size)


	uintSize := uint32(size)

	for i := 0; i < size; i++ {
		testRatings[i] = ratings[rand.Uint32() % uintSize]
	}

	return testRatings

}