package service

import (
	"github.com/gin-gonic/gin"
	"github.com/iliecirciumaru/rs-backend/model"
	"net/http"
	"github.com/iliecirciumaru/rs-backend/structs"
)

func NewRouteService(userS UserService, ratingS RatingService) RouteService {
	return RouteService{
		userService: userS,
		ratingService: ratingS,
	}
}

type RouteService struct {
	userService UserService
	ratingService RatingService
}

func (s *RouteService) RegisterUser(c *gin.Context) {
	request := model.UserRegisterRequest{}
	c.BindJSON(&request)

	err := s.userService.RegisterUser(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.CustomError{Message: err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "User successfully registered",
	})
}

func (s *RouteService) LoginUser(c *gin.Context) {
	request := model.UserLoginRequest{}
	c.BindJSON(&request)

	token, err := s.userService.Login(request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, structs.CustomError{Message: err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func (s *RouteService) AddRating(c *gin.Context) {
	user, _ := c.Get("user")
	u, _ := user.(model.User)

	request := model.RatingAddRequest{}
	c.BindJSON(&request)

	err := s.ratingService.AddRating(u, request)

	if err != nil {
		c.JSON(http.StatusBadRequest, structs.CustomError{Message: err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "rating successfully added",
	})
}