package model

import (
	"sync"
	"sort"
)

//func (r *Recommendation) CalculateMovieSimilarity(mID int64, uRats map[int64]float64, wg *sync.WaitGroup, mut *sync.Mutex, movieUserRatings map[int64]map[int64]float64) {
//	defer wg.Done()
//	cosineSimilarities := make([]Similarity, 0, len(movieUserRatings)-1)
//	for m2ID, u2Ratings := range movieUserRatings {
//		if mID == m2ID {
//			continue
//		}
//
//		cosineSimilarities = append(cosineSimilarities, Similarity{
//			ID:    m2ID,
//			Value: r.cosineSimilarity(uRats, u2Ratings),
//		})
//	}
//
//	sort.Sort(BySimilarityDesc(cosineSimilarities))
//	mut.Lock()
//	r.movieSimilarties[mID] = cosineSimilarities
//	mut.Unlock()
//}

func (r *Recommendation) CalculateMovieSimilarity(jobs chan int64, wg *sync.WaitGroup, mut *sync.Mutex, movieUserRatings map[int64]map[int64]float64) {
	var uRats map[int64]float64
	for mID := range jobs {
		uRats, _ = movieUserRatings[mID]

		cosineSimilarities := make([]Similarity, 0, len(movieUserRatings)-1)
		for m2ID, u2Ratings := range movieUserRatings {
			if mID == m2ID {
				continue
			}

			cosineSimilarities = append(cosineSimilarities, Similarity{
				ID:    m2ID,
				Value: r.cosineSimilarity(uRats, u2Ratings),
			})
		}

		sort.Sort(BySimilarityDesc(cosineSimilarities))
		mut.Lock()
		r.movieSimilarties[mID] = cosineSimilarities
		mut.Unlock()

		wg.Done()
	}
}