package main

import (
	"fmt"
	"math"
	"time"
	"sync"
)

func main() {
	t1 := time.Now().Unix()
	//r := model.Recommendation{}

	//fmt.Println(r.MeanRating([]float64{1.1, 2.2, 3.3, 4.4}))

	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fmt.Println(a[5:10])

	fmt.Println(math.Sqrt(0))

	t2 := time.Now().Unix()

	fmt.Println(t1 - t2)


	emulateCosineSimilarityCalculation()
}


var similarities map[int64]float64

func emulateCosineSimilarityCalculation() {
	var wg sync.WaitGroup
	mutex := sync.Mutex{}

	similarities = make(map[int64]float64)

	simils := []float64{1,2,3,4,5,6,7,8,9,10}
	for i, _ := range simils {
		wg.Add(1)
		go calculateSimilarity(i, simils, &wg, &mutex)
	}



	wg.Wait()
	fmt.Println("Calculation is done")
	fmt.Println(similarities)
}

func calculateSimilarity(id int, data []float64, wg *sync.WaitGroup, mutex *sync.Mutex) {
	mutex.Lock()
	similarities[int64(id)] = data[id]
	mutex.Unlock()
	wg.Done()
}


