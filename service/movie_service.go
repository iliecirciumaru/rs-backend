package service

import (
	"github.com/iliecirciumaru/rs-backend/repo"
	"github.com/iliecirciumaru/rs-backend/model"
)

func NewMovieService(movieRepo repo.MovieRepo, ratingRepo repo.RatingRepo) MovieService {
	return MovieService{
		repo:movieRepo,
		ratingRepo: ratingRepo,
	}
}

type MovieService struct {
	repo repo.MovieRepo
	ratingRepo repo.RatingRepo
}

func (s *MovieService) GetMovieWithUserRating(movieID int64, user model.User) (model.MovieView, error) {
	var movieView model.MovieView

	movie, err := s.repo.GetMovie(movieID)
	if err != nil {
		return movieView, err
	}


	movieView.ID = movieID
	movieView.Information = movie.Information
	movieView.Title = movie.Title


	rating := s.ratingRepo.GetRatingByMovieUserID(user.ID, movieID)
	movieView.UserRating = rating.Value

	return movieView, nil
}