package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at"`
}

func (u *User) GetBSON() (interface{}, error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	type my *User
	return my(u), nil
}
