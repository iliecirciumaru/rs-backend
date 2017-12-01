package service

import (
	"github.com/gin-gonic/gin"
	"github.com/iliecirciumaru/rs-backend/model"
	"net/http"
	"github.com/iliecirciumaru/rs-backend/structs"
)

func NewRouteService(userService UserService) RouteService {
	return RouteService{
		userService: userService,
	}
}

type RouteService struct {
	userService UserService
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