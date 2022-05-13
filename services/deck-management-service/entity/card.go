package entity

import "time"

type Card struct {
	ID        string     `bson:"_id,omitempty"`
	Question  string     `bson:"question"`
	Answer    string     `bson:"answer"`
	UserID    string     `bson:"user_id"`
	DeckID    string     `bson:"deck_id"`
	CreatedAt *time.Time `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

type CardReq struct {
	ID       string `json:"id,omitempty"`
	Question string `json:"question"     binding:"required"`
	Answer   string `json:"answer"       binding:"required"`
	DeckID   string `json:"deckID"       binding:"required"`
}

type CardRes struct {
	ID       string `json:"id"`
	Question string `json:"question" binding:"required"`
	Answer   string `json:"answer"   binding:"required"`
	DeckID   string `json:"deckID"   binding:"required"`
}

type CardUseCaseInterface interface {
	CreateCard(deckID, userID string, card *CardReq) (*CardRes, error)
	CreateCards(deckID, userID string, card []CardReq) ([]CardRes, error)
	UpdateCard(cardID, userID, deckID string, card *CardReq) (*CardRes, error)
	DeleteCard(userID, deckID, cardID string) error
}

type CardStoreInterface interface {
	SaveCard(deckID, userID string, card *Card) (string, error)
	SaveCards(deckID, userID string, card []Card) ([]string, error)
	UpdateCard(cardID, userID, deckID string, card *Card) error
	DeleteCard(userID, deckID, cardID string) error
}
