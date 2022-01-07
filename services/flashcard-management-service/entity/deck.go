package entity

import "time"

type Deck struct {
	ID        string    `json:"id,omitempty" bson:"id"`
	Name      string    `json:"name"         bson:"name"       validate:"required"`
	UserID    string    `json:"-"            bson:"user_id"`
	CreatedAt time.Time `json:"-"            bson:"created_at"`
	UpdatedAt time.Time `json:"-"            bson:"updated_at"`
	DeletedAt time.Time `json:"-"            bson:"deleted_at"`
}
