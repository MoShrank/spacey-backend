package entities

import "time"

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name"         bson:"name"`
	Email     string    `json:"email"        bson:"email"`
	Password  string    `json:"password"     bson:"password"`
	createdAt time.Time `                    bson:"created_at,omitempty"`
	updatedAt time.Time `                    bson:"updated_at,omitempty"`
	deletedAt time.Time `                    bson:"deleted_at,omitempty"`
}
