package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iliecirciumaru/rs-backend/db"
	"log"
	"github.com/iliecirciumaru/rs-backend/service"
	"github.com/iliecirciumaru/rs-backend/repo"
	"upper.io/db.v3/lib/sqlbuilder"
	"github.com/iliecirciumaru/rs-backend/middleware"
)

func main() {
	r := gin.Default()

	dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	if err != nil {
		log.Fatal(err)
	}

	userRepo, ratingRepo := initRepos(dbsess)
	routes, _:= initServices(userRepo, ratingRepo)

	r.Use(middleware.AuthValidation(userRepo))

	bootstrapRoutes(r, routes)


	r.Run()
}

func bootstrapRoutes(r *gin.Engine, routes service.RouteService) {
	r.POST("/user", routes.RegisterUser)
	r.POST("/login", routes.LoginUser)
	r.POST("/rating", routes.AddRating)
}

func initServices(userRepo repo.UserRepo, ratingRepo repo.RatingRepo) ( service.RouteService, service.UserService) {
	userService := service.NewUserService(userRepo)
	ratingService := service.NewRatingService(ratingRepo)
	routes := service.NewRouteService(userService, ratingService)

	return routes, userService
}

func initRepos(db sqlbuilder.Database) (repo.UserRepo, repo.RatingRepo) {
	userRepo := repo.NewUserRepo(db)
	ratingRepo := repo.NewRatingRepo(db)
	return userRepo, ratingRepo
}

