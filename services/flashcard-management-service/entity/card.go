package entity

import "time"

type Card struct {
	ID        string     `bson:"_id"`
	Question  string     `bson:"question"`
	Answer    string     `bson:"answer"`
	DeckID    string     `bson:"deck_id"`
	UserID    string     `bson:"user_id"`
	CreatedAt *time.Time `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

type CardReq struct {
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer"   validate:"required"`
	DeckID   string `json:"deck_id"  validate:"required"`
}

type CardRes struct {
	ID string `json:"id"`
	CardReq
}

type CardUseCaseInterface interface {
	CreateCard(userID string, card *CardReq) (*CardRes, error)
	GetCardsByDeckID(userID, deckID string) ([]CardRes, error)
	GetCardByID(userID, cardID string) (*CardRes, error)
	UpdateCard(userID, cardID string, card *CardReq) (*CardRes, error)
	DeleteCard(userID, cardID string) error
}

type CardStoreInterface interface {
	SaveCard(card *Card) (string, error)
	GetCardsByDeckID(userID, deckID string) ([]Card, error)
	GetCardByID(userID, cardID string) (*Card, error)
	UpdateCard(userID string, card *Card) error
	DeleteCard(userID, cardID string) error
}
