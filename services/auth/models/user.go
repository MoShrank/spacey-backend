package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	CreatedAtTs time.Time          `json:"-" bson:"created_at_ts"`
	UpdatedAtTs time.Time          `json:"-" bson:"updated_at_ts"`
	DeletedAtTs time.Time          `json:"-" bson:"deleted_at_ts"`
}
