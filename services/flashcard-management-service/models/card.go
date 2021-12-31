package models

import "time"

type Card struct {
	ID       string `bson:"_id"`
	Question string `bson:"question"   json:"question" validate:"required"`
	Answer   string `bson:"answer"     json:"answer"   validate:"required"`
	DeckID   string `bson:"deck_id"    json:"deck_id"  validate:"required"`
	//Deck      Deck   `bson:"deck"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}
