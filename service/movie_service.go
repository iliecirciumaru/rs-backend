package service

import (
	"github.com/iliecirciumaru/rs-backend/model"
	"github.com/iliecirciumaru/rs-backend/repo"
	"sort"
)

func NewMovieService(movieRepo repo.MovieRepo, ratingRepo repo.RatingRepo, recommendNeighbours uint) MovieService {

	return MovieService{
		repo:       movieRepo,
		ratingRepo: ratingRepo,
		rec:        model.Recommendation{UuNeighbours: recommendNeighbours},
	}
}

type MovieService struct {
	repo       repo.MovieRepo
	ratingRepo repo.RatingRepo
	rec        model.Recommendation
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
	if movie.PosterURL != nil {
		movieView.PosterURL = *movie.PosterURL
	}

	rating := s.ratingRepo.GetRatingByMovieUserID(user.ID, movieID)
	// TODO if user hasn't scored yet, try to predict
	if rating.Value == 0 {

	} else {
		movieView.Rated = true
	}

	avg, err := s.ratingRepo.GetAVGMovieRating([]int64{movieID})
	if err == nil && len(avg) == 1 {
		movieView.AverageRating = model.RoundRating(avg[0].Value)
	}

	movieView.UserRating = rating.Value
	return movieView, nil
}

func (s *MovieService) GetTopRatedMovies(number int) ([]model.MovieView, error) {
	ratings, err := s.ratingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	filmAverages := s.rec.ComputeFilmAverageScores(ratings)

	movieViews := make([]model.MovieView, len(filmAverages))
	i := 0
	for movieID, mean := range filmAverages {
		movieViews[i].ID = movieID

		// soft round till 2 numbers
		movieViews[i].AverageRating = model.RoundRating(mean)
		i++
	}

	sort.Slice(movieViews, func(i, j int) bool {
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
				if movie.PosterURL != nil {
					movieViews[i].PosterURL = *movie.PosterURL
				}
				break
			}
		}
	}

	return movieViews[:number], nil
}

func (s *MovieService) GetRecommendationForUser(user model.User, number int) ([]model.UserMovieView, error) {
	ratings, err := s.ratingRepo.GetAll()
	if err != nil {
		return nil, err
	}
	predictions := s.rec.PredictUserScoreUUCLF(user.ID, ratings)
	if len(predictions) < number {
		number = len(predictions)
	}

	result := make([]model.UserMovieView, number)
	if number == 0 {
		return result, nil
	}

	predictions = predictions[:number]
	movieIDs := make([]int64, len(predictions))
	for i, p := range predictions {
		movieIDs[i] = p.MovieID
	}

	movies, err := s.repo.GetMovieByIDs(movieIDs)
	if err != nil {
		return nil, err
	}

	for i, movie := range movies {
		result[i] = model.UserMovieView{
			ID:          movie.ID,
			Information: movie.Information,
			Title:       movie.Title,
			Rated:       false,
		}

		if movie.PosterURL != nil {
			result[i].PosterURL = *movie.PosterURL
		}

		for _, p := range predictions {
			if p.MovieID == movie.ID {
				result[i].UserRating = model.RoundRating(p.PredictedScore)
				break
			}
		}
	}

	return result, nil
}

func (s *MovieService) GetRecentReleases(number int) ([]model.MovieView, error) {
	movies, err := s.repo.GetLatestMovies(number)
	if err != nil {
		return nil, err
	}

	result := s.ConstructMovieViews(movies)


	return result, nil
}

func (s *MovieService) GetMovieByPrefix(prefix string) ([]model.MovieView, error) {
	movies, err := s.repo.GetMoviesByPrefix(prefix)
	if err != nil {
		return nil, err
	}

	result := s.ConstructMovieViews(movies)
	return result, nil
}

func (s *MovieService) ConstructMovieViews(movies []model.Movie) []model.MovieView {
	result := make([]model.MovieView, len(movies))

	if len(movies) == 0 {
		return result
	}

	movieIDs := make([]int64, len(movies))


	for i, m := range movies {
		movieIDs[i] = m.ID

		result[i] = model.MovieView{
			ID:          m.ID,
			Information: m.Information,
			Title:       m.Title,
		}

		if m.PosterURL != nil {
			result[i].PosterURL = *m.PosterURL
		}
	}

	avgs, err := s.ratingRepo.GetAVGMovieRating(movieIDs)
	if err == nil {
		for _, avg := range avgs {
			for i := 0; i < len(result); i++ {
				if result[i].ID == avg.Key {
					result[i].AverageRating = model.RoundRating(avg.Value)
					break
				}
			}
		}
	}

	return result
}

func (s *MovieService) GetSimilarMovies(movieID int64, number int) ([]model.MovieView, error) {
	ratings, err := s.ratingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	similarities := s.rec.GetMostSimilarMovies(movieID, ratings)
	if len(similarities) == 0 {
		return s.GetTopRatedMovies(number)
	}

	if len(similarities) < number {
		number = len(similarities)
	}

	similarities = similarities[:number]
	movieIDs := make([]int64, number)
	for i, sim := range similarities {
		movieIDs[i] = sim.ID
	}

	movies, err := s.repo.GetMovieByIDs(movieIDs)
	if err != nil {
		return nil, err
	}

	result := s.ConstructMovieViews(movies)


	return result, nil
}