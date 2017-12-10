package service

import (
	"fmt"
	"github.com/iliecirciumaru/rs-backend/model"
	"github.com/iliecirciumaru/rs-backend/repo"
)

func NewUserService(userRepo repo.UserRepo) UserService {
	return UserService{repo: userRepo}
}

type UserService struct {
	repo repo.UserRepo
}

func (s *UserService) RegisterUser(request model.UserRegisterRequest) error {
	user := model.User{Login: request.Login, Password: request.Password, Name: request.Name}

	return s.repo.AddUser(user)
}

func (s *UserService) Login(request model.UserLoginRequest) (model.UserLoginResponse, error) {
	user, err := s.repo.GetUserByLoginAndPassword(request.Login, request.Password)

	if err != nil {
		return model.UserLoginResponse{}, fmt.Errorf("Login and Password pair is invalid")
	}

	return model.UserLoginResponse{
		Token: user.Login + user.Name,
		Login: user.Login,
		Name:  user.Name,
	}, nil
}
