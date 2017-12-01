package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iliecirciumaru/rs-backend/db"
	"log"
	"github.com/iliecirciumaru/rs-backend/service"
	"github.com/iliecirciumaru/rs-backend/repo"
)


func main() {
	r := gin.Default()

	bootstrapRoutes(r)

	r.Run()
}

func bootstrapRoutes(r *gin.Engine) {
	routes := initDependincies()
	r.POST("/user", routes.RegisterUser)
}

func initDependincies() service.RouteService{
	dbsess, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repo.NewUserRepo(dbsess)
	userService := service.NewUserService(userRepo)
	routes := service.NewRouteService(userService)

	return routes
}

