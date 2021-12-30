package models

type Card struct {
	ID       string `bson:"_id"`
	Question string `bson:"question"   json:"question" validate:"required"`
	Answer   string `bson:"answer"     json:"answer"   validate:"required"`
	DeckID   string `bson:"deck_id"    json:"deck_id"  validate:"required"`
	//Deck      Deck   `bson:"deck"`
	UserID    string `bson:"user_id"`
	CreatedAt string `bson:"created_at"`
	UpdatedAt string `bson:"updated_at"`
	DeletedAt string `bson:"deleted_at"`
}
