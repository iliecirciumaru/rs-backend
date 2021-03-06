package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iliecirciumaru/rs-backend/db"
	"github.com/iliecirciumaru/rs-backend/middleware"
	"github.com/iliecirciumaru/rs-backend/repo"
	"github.com/iliecirciumaru/rs-backend/service"
	"log"
	"upper.io/db.v3/lib/sqlbuilder"
)

func main() {
	r := gin.Default()

	dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	if err != nil {
		log.Fatal(err)
	}

	userRepo, ratingRepo, movieRepo := initRepos(dbsess)
	routes, _, _ := initServices(userRepo, ratingRepo, movieRepo)

	r.Use(middleware.AuthValidation(userRepo))
	r.Use(middleware.CORSMiddleware())

	bootstrapRoutes(r, routes)

	r.Run()
}

func bootstrapRoutes(r *gin.Engine, routes service.RouteService) {
	r.POST("/user", routes.RegisterUser)
	r.POST("/login", routes.LoginUser)
	r.POST("/rating", routes.AddRating)
	r.GET("/movie/:id", routes.GetMovie)
	r.GET("/movie/:id/similar", routes.GetSimilar)
	r.GET("/movies/top", routes.GetTopMovies)
	r.GET("/movies/recommended", routes.GetReccommendedMovies)
	r.GET("/movies/recent", routes.GetRecentReleases)
	r.GET("/movies/search", routes.GetSearchMovies)
}

func initServices(userRepo repo.UserRepo, ratingRepo repo.RatingRepo, movieRepo repo.MovieRepo) (
	service.RouteService, service.UserService, service.MovieService) {
	userService := service.NewUserService(userRepo)
	ratingService := service.NewRatingService(ratingRepo)
	movieService := service.NewMovieService(movieRepo, ratingRepo, 8)
	routes := service.NewRouteService(userService, ratingService, movieService)

	return routes, userService, movieService
}

func initRepos(db sqlbuilder.Database) (repo.UserRepo, repo.RatingRepo, repo.MovieRepo) {
	userRepo := repo.NewUserRepo(db)
	ratingRepo := repo.NewRatingRepo(db)
	movieRepo := repo.NewMovieRepo(db)
	return userRepo, ratingRepo, movieRepo
}
