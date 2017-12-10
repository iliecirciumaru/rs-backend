package repo

import (
	"github.com/iliecirciumaru/rs-backend/model"
	"upper.io/db.v3/lib/sqlbuilder"
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
	_, err := r.db.Collection("users").Insert(&user)

	return err
}

func (r *UserRepo) ValidateUserByToken(token string) (model.User, error) {
	query := r.db.SelectFrom("users").Where("CONCAT(login, name) = ?", token)

	var user model.User
	err := query.One(&user)
	return user, err
}

func (r *UserRepo) GetUserByLoginAndPassword(login, password string) (model.User, error) {
	query := r.db.SelectFrom("users").Where("login = ? AND password = ?", login, password)

	var user model.User
	err := query.One(&user)
	return user, err
}
