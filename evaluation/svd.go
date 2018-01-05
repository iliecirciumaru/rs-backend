package main

import (
	"gonum.org/v1/gonum/mat"
	"log"
	"github.com/iliecirciumaru/rs-backend/db"
	"fmt"
	"github.com/iliecirciumaru/rs-backend/repo"
	"upper.io/db.v3/lib/sqlbuilder"
	"github.com/iliecirciumaru/rs-backend/model"
	"sort"
)

func main() {
	//data := make([]float64, 36)
	//for i := range data {
	//	data[i] = rand.NormFloat64()
	//}
	//a := mat.NewDense(6, 6, data)



	dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	if err != nil {
		log.Fatal(err)
	}

	u, v, singValues := computeSVD(computeMatrix(dbsess))

	fmt.Println("Users matrix")
	fmt.Println( mat.Formatted(u, mat.Prefix("\t    "), mat.Excerpt(3)))

	fmt.Println("Movie matrix")
	fmt.Println( mat.Formatted(v, mat.Prefix("\t    "), mat.Excerpt(3)))

	fmt.Println("Singular Values")
	fmt.Printf("Length: %v\n, %v...\n", len(singValues), singValues[0:3])


	uID := uint(1)
	mID := uint(1)

	score := predictScore(uID, mID, u, v, singValues)

	fmt.Printf("Predicted for User %v, Movie %v - %5.2f\n", uID, mID, score)
}

func computeMatrix(dbsess sqlbuilder.Database) *mat.Dense {
	ratingRepo := repo.NewRatingRepo(dbsess)
	movieRepo := repo.NewMovieRepo(dbsess)
	userRepo := repo.NewUserRepo(dbsess)

	rats, err := ratingRepo.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	movies, err := movieRepo.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ID < movies[j].ID
	})

	users, err := userRepo.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(users)

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})


	start := 0
	end := int(float64(len(rats)) * 0.8)
	testRats := rats[start:end]


	var recommender model.Recommendation = model.Recommendation{UuNeighbours:15}
	userMovieRatings := recommender.GetUserMovieRatings(testRats)


	uCount := len(users)
	mCount := len(movies)

	var movieRatings map[int64]float64
	var ok bool
	var r float64

	matrix := make([]float64, uCount * mCount, uCount * mCount)
	for i, user := range users {
		if movieRatings, ok = userMovieRatings[user.ID]; ok {
			for j, movie := range movies {
				if r, ok = movieRatings[movie.ID]; ok {
					matrix[i * mCount + j] = r
				}
			}
		}
	}


	a := mat.NewDense(uCount, mCount, matrix)

	fmt.Println("Matrix computation finished")
	fmt.Println( mat.Formatted(a, mat.Prefix("\t    "), mat.Excerpt(10)))

	return a
}

func computeSVD(a *mat.Dense) (u *mat.Dense, v *mat.Dense, singValues []float64) {
	svd := mat.SVD{}

	ok := svd.Factorize(a, mat.SVDFull)
	if !ok {
		log.Fatalf("SVD factorization failed")
	}

	u = svd.UTo(nil)
	v = svd.VTo(nil)
	singValues = svd.Values(nil)

	return
}

func predictScore(userID, movieID uint, u, v *mat.Dense, singValues []float64) float64 {
	ufeatures := u.RowView(2)


	//for movie := 0; movie < 200; movie++ {
	//	mfeatures := v.RawRowView(movie)
	//	//mfeatures := v.ColView(movie)
	//
	//	score := float64(0)
	//	n := len(ufeatures)
	//	//if n > len(mfeatures) {
	//	//	n = len(mfeatures)
	//	//}
	//
	//
	//	for i := 0; i < n; i++ {
	//		//if i == 1 {
	//		//	score += ufeatures[i] * mfeatures[i] * singValues[i]
	//		//} else {
	//			score += ufeatures[i] * mfeatures[i]
	//		//}
	//
	//	}
	//	fmt.Printf("%v. %5.2f\n", movie, score)
	//}

	for movie := 0; movie < 200; movie++ {
		fmt.Printf("%v. %5.2f\n", movie, mat.Dot(ufeatures, v.RowView(movie)))
		fmt.Printf("%v. %5.2f\n", movie, mat.Dot(ufeatures, v.ColView(movie)))
	}



	return 4.23
}