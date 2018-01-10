package model

import (
	"sync"
	"sort"
)

//func (r *Recommendation) CalculateMovieSimilarity(mID int64, uRats map[int64]float32, wg *sync.WaitGroup, mut *sync.Mutex, movieUserRatings map[int64]map[int64]float32) {
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
//	r.MovieSimilarties[mID] = cosineSimilarities
//	mut.Unlock()
//}

func (r *Recommendation) CalculateMovieSimilarity(jobs chan int64, wg *sync.WaitGroup, mut *sync.Mutex, movieUserRatings map[int64]map[int64]float32) {
	var uRats map[int64]float32
	var simils2 []Similarity
	var ok bool
	for mID := range jobs {
		uRats, _ = movieUserRatings[mID]

		cosineSimilarities := make([]Similarity, 0, len(movieUserRatings)-1)
		for m2ID, u2Ratings := range movieUserRatings {
			if mID == m2ID {
				continue
			}

			mut.Lock()
			simils2, ok = r.MovieSimilarties[m2ID];
			mut.Unlock()

			if  ok {
				for _, s := range simils2 {
					if s.ID == mID {
						cosineSimilarities = append(cosineSimilarities, Similarity{
							ID:    m2ID,
							Value: s.Value,
						})
						break
					}
				}
			} else {
				cosineSimilarities = append(cosineSimilarities, Similarity{
					ID:    m2ID,
					Value: r.cosineSimilarity(uRats, u2Ratings),
				})
			}




		}

		sort.Sort(BySimilarityDesc(cosineSimilarities))
		mut.Lock()
		r.MovieSimilarties[mID] = cosineSimilarities
		mut.Unlock()

		wg.Done()
	}
}