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

	fmt.Println("USER SERVICE")
	fmt.Println(user)

	return s.repo.AddUser(user)
}
