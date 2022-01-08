package entity

import "time"

type Deck struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `                  jbson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

type DeckReq struct {
	Name string `json:"name" validate:"required"`
}

type DeckRes struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

type DeckUseCaseInterface interface {
	CreateDeck(userID string, deck *DeckReq) (*DeckRes, error)
	GetDecks(userID string) ([]DeckRes, error)
	GetDeck(userID, DeckID string) (*DeckRes, error)
	UpdateDeck(userID, DeckID string, deck *DeckReq) (*DeckRes, error)
	DeleteDeck(userID, DeckID string) error
}

type DeckStoreInterface interface {
	Save(deck *Deck) error
	FindAll(userID string) ([]Deck, error)
	Find(userID, deckID string) (*Deck, error)
	Update(deck *Deck) error
	Delete(userID, deckID string) error
}
