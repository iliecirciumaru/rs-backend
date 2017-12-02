package service

import (
	"github.com/iliecirciumaru/rs-backend/repo"
	"github.com/iliecirciumaru/rs-backend/model"
	"fmt"
)

func NewUserService(userRepo repo.UserRepo) UserService {
	return UserService{repo:userRepo}
}

type UserService struct {
	repo repo.UserRepo
}

func (s *UserService) RegisterUser(request model.UserRegisterRequest) error {
	user := model.User{Login: request.Login, Password:request.Password, Name: request.Name}

	return s.repo.AddUser(user)
}

func (s *UserService) Login(request model.UserLoginRequest) (string, error) {
	user, err := s.repo.GetUserByLoginAndPassword(request.Login, request.Password)

	if err != nil {
		return "", fmt.Errorf("Login and Password pair is invalid")
	}

	return user.Login + user.Name, nil

}
