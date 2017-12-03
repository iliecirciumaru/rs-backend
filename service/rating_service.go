package service

import (
	"fmt"
	"github.com/iliecirciumaru/rs-backend/model"
	"github.com/iliecirciumaru/rs-backend/repo"
	"time"
)

func NewRatingService(ratingRepo repo.RatingRepo) RatingService {
	return RatingService{
		repo: ratingRepo,
	}
}

type RatingService struct {
	repo repo.RatingRepo
}

func (s *RatingService) AddRating(user model.User, request model.RatingAddRequest) error {
	if request.Value < 0.1 || request.Value > 5.1 {
		return fmt.Errorf("Invalid rating value, should be beetwen 0.5 and 5, got %f", request.Value)
	}

	// TODO check that movie exists
	ts := time.Now().Unix()

	rating := model.Rating{
		Value:     request.Value,
		MovieID:   request.MovieID,
		UserID:    user.ID,
		Timestamp: ts,
	}

	return s.repo.AddRating(rating)
}
