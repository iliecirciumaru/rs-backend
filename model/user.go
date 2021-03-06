package model

type User struct {
	ID       int64  `json:"id" db:"id,omitempty"`
	Login    string `json:"login" db:"login"`
	Password string `json:"-" db:"password"`
	Name     string `json:"name" db:"name"`
}

type UserRegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
