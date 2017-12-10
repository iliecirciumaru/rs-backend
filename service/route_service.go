package service

import (
	"github.com/gin-gonic/gin"
	"github.com/iliecirciumaru/rs-backend/model"
	"github.com/iliecirciumaru/rs-backend/structs"
	"net/http"
	"strconv"
)

func NewRouteService(userS UserService, ratingS RatingService, movieS MovieService) RouteService {
	return RouteService{
		userService:   userS,
		ratingService: ratingS,
		movieService:  movieS,
	}
}

type RouteService struct {
	userService   UserService
	ratingService RatingService
	movieService  MovieService
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

	loginReponse, err := s.userService.Login(request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, structs.CustomError{Message: err.Error()})
		return
	}

	c.JSON(200, loginReponse)
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

func (s *RouteService) GetMovie(c *gin.Context) {
	user, _ := c.Get("user")
	u, _ := user.(model.User)

	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.CustomError{Message: "MovieID should be numeric"})
		return
	}

	movieView, err := s.movieService.GetMovieWithUserRating(int64(movieID), u)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.CustomError{Message: err.Error()})
		return
	}

	c.JSON(200, movieView)
}

func (s *RouteService) GetTopMovies(c *gin.Context) {
	snumber, _ := c.GetQuery("number")
	number, err := strconv.Atoi(snumber)
	if err != nil {
		number = 10
	}

	movieViews, err := s.movieService.GetTopRatedMovies(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.CustomError{Message: err.Error()})
		return
	}

	c.JSON(200, movieViews)
}

func (s *RouteService) GetReccommendedMovies(c *gin.Context) {
	snumber, _ := c.GetQuery("number")
	number, err := strconv.Atoi(snumber)
	if err != nil {
		number = 10
	}

	user, _ := c.Get("user")
	u, _ := user.(model.User)

	movieViews, err := s.movieService.GetRecommendationForUser(u, number)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.CustomError{Message: err.Error()})
		return
	}

	c.JSON(200, movieViews)
}

func (s *RouteService) GetRecentReleases(c *gin.Context) {
	snumber, _ := c.GetQuery("number")
	number, err := strconv.Atoi(snumber)
	if err != nil {
		number = 10
	}

	movieViews, err := s.movieService.GetRecentReleases(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.CustomError{Message: err.Error()})
		return
	}

	c.JSON(200, movieViews)
}
