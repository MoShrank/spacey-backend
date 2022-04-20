package entity

import (
	"errors"
	"time"
)

type User struct {
	ID          string     `bson:"_id,omitempty"`
	Name        string     `bson:"name,omitempty"`
	Email       string     `bson:"email"`
	Password    string     `bson:"password"`
	BetaUser    bool       `bson:"betaUser"`
	CreatedAtTs *time.Time `bson:"created_at_ts"`
	UpdatedAtTs *time.Time `bson:"updated_at_ts"`
	DeletedAtTs *time.Time `bson:"deleted_at_ts"`
}

type UserResponseModel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	BetaUser bool   `json:"betaUser"`
	Token    string `json:"token,omitempty"`
}

type UserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserStoreInterface interface {
	SaveUser(user *User) (string, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
}

type UserUsecaseInterface interface {
	CreateUser(user interface{}) (*UserResponseModel, error)
	Login(email, password string) (*UserResponseModel, error)
	GetUserByID(id string) (*UserResponseModel, error)
}

var ErrEmailAlreadyExists = errors.New("email already exists.")
