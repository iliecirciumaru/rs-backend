package model

import (
	"sort"
	"fmt"
	"time"
	"github.com/iliecirciumaru/rs-backend/structs"
)

type ClusteringUtility struct {
	Rec Recommendation
}

func (c *ClusteringUtility) ExtractRatings(ratings []Rating, movieSet []int64) []Rating {
	newRats := make([]Rating, 0, 256)
	for _, r := range ratings {
		if in(r.MovieID, movieSet) {
			newRats = append(newRats, r)
		}
	}

	return newRats
}

func (c *ClusteringUtility) Cluster(ratings []Rating, mostRatedMovie int64) map[int64][]int64{
	fmt.Println("Start clustering")
	start := time.Now().UnixNano()
	movieUserRatings := c.Rec.getMovieUserRatings(ratings)
	c.Rec.normalizeUserMovieOrMovieUserRatings(movieUserRatings)

	//clustersNum := 3

	mostRatedMovieRats, _ := movieUserRatings[mostRatedMovie]

	similarities := make([]Similarity, len(movieUserRatings) - 1)

	i := 0
	for mID, uRats := range movieUserRatings {
		if mID == mostRatedMovie {
			continue
		}

		similarities[i] = Similarity{
			ID: mID,
			Value: c.Rec.cosineSimilarity(mostRatedMovieRats, uRats),
		}

		i++
	}

	sort.Sort(BySimilarityDesc(similarities))


	centroidTwo := similarities[len(similarities) / 2]
	centroidThree := similarities[len(similarities) - 1]

	//fmt.Println(similarities[0:5])

	centroids := []int64{mostRatedMovie, centroidTwo.ID, centroidThree.ID}
	cluster := make(map[int64][]int64)
	var maxSimilarity, tempSimilarity float64
	//var maxSimilarity float64
	var centroid int64
	//resultChannel := make(chan structs.KeyValue)
	//var similarToCentroid structs.KeyValue


	for _, centroid = range centroids {
		cluster[centroid] = make([]int64, 0, 64)
	}

	for mID, uRats := range movieUserRatings {
		if in(mID, centroids) {
			cluster[mID] = append(cluster[mID], mID)
			continue
		}
		maxSimilarity = -2000

		for _, centr := range centroids {
			tempSimilarity = c.Rec.cosineSimilarity(movieUserRatings[centr], uRats)
			if tempSimilarity > maxSimilarity {
				maxSimilarity = tempSimilarity
				centroid = centr
			}
		}

		//for _, centr := range centroids {
		//	go c.CosineSimilarityAsync(centr, movieUserRatings[centr], uRats, resultChannel)
		//}
		//
		//for i := 0; i < clustersNum; i++ {
		//	similarToCentroid = <- resultChannel
		//	if similarToCentroid.Value > maxSimilarity {
		//		maxSimilarity = similarToCentroid.Value
		//		centroid = similarToCentroid.Key
		//	}
		//}

		cluster[centroid] = append(cluster[centroid], mID)
	}

	end := time.Now().UnixNano()
	fmt.Printf("Clustering is finished, time: %.2fs\n", float64(end-start) / 1000000000)


	//fmt.Printf("Centroid 1: %v, rates: %v\n", mostRatedMovie, len(mostRatedMovieRats))
	//fmt.Printf("Centroid 2: %v, rates: %v, similarity: %.3f\n",
	//	centroidTwo.ID, len(movieUserRatings[centroidTwo.ID]), centroidTwo.Value)
	//fmt.Printf("Centroid 3: %v, rates: %v, similarity: %.3f\n",
	//	centroidThree.ID, len(movieUserRatings[centroidThree.ID]), centroidThree.Value)

	for centroid, neigh := range cluster {
		fmt.Printf("Centroid %v, neigh %v, rates %v\n", centroid, len(neigh), len(movieUserRatings[centroid]))
	}

	return cluster
}


func (c *ClusteringUtility) CosineSimilarityAsync(centroid int64, ratings1, ratings2 map[int64]float64, result chan<- structs.KeyValue) {
	similarity := c.Rec.cosineSimilarity(ratings1, ratings2)
	result <- structs.KeyValue{Key: centroid, Value: similarity}
}

func in(el int64, a []int64) bool {
	for _, e := range a {
		if e == el {
			return true
		}
	}

	return false
}

