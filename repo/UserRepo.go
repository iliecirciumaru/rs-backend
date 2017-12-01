package repo

import (
	"github.com/iliecirciumaru/rs-backend/model"
	"upper.io/db.v3/lib/sqlbuilder"
	"fmt"
)

func NewUserRepo(db sqlbuilder.Database) UserRepo {
	return UserRepo{
		db: db,
	}
}

type UserRepo struct {
	db sqlbuilder.Database
}

func (r *UserRepo) AddUser(user model.User) error {
	fmt.Println("USER REPO")
	fmt.Println(user)

	_, err := r.db.Collection("users").Insert(&user)

	return err
}
