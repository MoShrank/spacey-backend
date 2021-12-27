package models

import (
	"time"
)

type User struct {
	ID          string    `json:"-"`
	Name        string    `json:"name"     bson:"name"          validate:"required"`
	Email       string    `json:"email"    bson:"email"         validate:"required,email"`
	Password    string    `json:"password" bson:"password"      validate:"required"`
	CreatedAtTs time.Time `json:"-"        bson:"created_at_ts"`
	UpdatedAtTs time.Time `json:"-"        bson:"updated_at_ts"`
	DeletedAtTs time.Time `json:"-"        bson:"deleted_at_ts"`
}
