package service

import (
	"github.com/iliecirciumaru/rs-backend/model"
	"github.com/iliecirciumaru/rs-backend/repo"
	"sort"
)

func NewMovieService(movieRepo repo.MovieRepo, ratingRepo repo.RatingRepo) MovieService {

	return MovieService{
		repo:       movieRepo,
		ratingRepo: ratingRepo,
		rec: model.Recommendation{},
	}
}

type MovieService struct {
	repo       repo.MovieRepo
	ratingRepo repo.RatingRepo
	rec model.Recommendation
}

func (s *MovieService) GetMovieWithUserRating(movieID int64, user model.User) (model.UserMovieView, error) {
	var movieView model.UserMovieView

	movie, err := s.repo.GetMovieByID(movieID)
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

func (s *MovieService) GetTopRatedMovies(number int) ([]model.MovieView, error) {
	ratings, err := s.ratingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	filmAverages := s.rec.ComputeFilmAverages(ratings)

	movieViews := make([]model.MovieView, len(filmAverages))
	i := 0
	for movieID, mean := range filmAverages {
		movieViews[i].ID = movieID
		movieViews[i].AverageRating = float64(int64(mean*20+0.5)) / 20
		i++
	}

	sort.Slice(movieViews, func(i,j int) bool {
		return movieViews[i].AverageRating > movieViews[j].AverageRating
	})

	movieViews = movieViews[:number]
	movieIDs := make([]int64, len(movieViews))
	for i, view := range movieViews {
		movieIDs[i] = view.ID
	}


	movies, _ := s.repo.GetMovieByIDs(movieIDs)

	for i, view := range movieViews {
		for _, movie := range movies {
			if view.ID == movie.ID {
				movieViews[i].Title = movie.Title
				movieViews[i].Information = movie.Information
				break
			}
		}
	}



	return movieViews[:number], nil
}
