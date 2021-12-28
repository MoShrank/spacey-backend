package models

type Flashcard struct {
	ID        string `bson:"_id"`
	Question  string `bson:"question"`
	Answer    string `bson:"answer"`
	DeckID    string `bson:"deck_id"`
	Deck      Deck   `bson:"deck"`
	UserID    string `bson:"user_id"`
	CreatedAt string `bson:"created_at"`
	UpdatedAt string `bson:"updated_at"`
	DeletedAt string `bson:"deleted_at"`
}
