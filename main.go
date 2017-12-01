package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iliecirciumaru/rs-backend/db"
	"log"
)


func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	_, err := db.GetDb("root", "password", "rs")
	if err != nil {
		log.Fatal(err)
	}

	r.Run()
}