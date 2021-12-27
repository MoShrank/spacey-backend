package store

import (
	"context"

	"github.com/moshrank/spacey-backend/pkg/logger"

	"github.com/moshrank/spacey-backend/services/flashcard-management-service/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeckStore struct {
	db     *mongo.Database
	logger logger.LoggerInterface
}

type DeckStoreInterface interface {
	GetDeck(id string) (*entities.Deck, error)
	GetDecks() ([]entities.Deck, error)
	CreateDeck(deck *entities.Deck) error
	UpdateDeck(deck *entities.Deck) error
	DeleteDeck(id string) error
}

func (store *DeckStore) GetDeck(userID string, id string) (*entities.Deck, error) {
	ctx := context.TODO()

	var deck entities.Deck
	err := store.db.Collection("Deck").
		FindOne(ctx, bson.M{"_id": id, "UserID": userID}).
		Decode(&deck)

	if err != nil {
		store.logger.Fatal(err)
	}

	return &deck, nil
}

func (store *DeckStore) GetDecks(userID string) ([]entities.Deck, error) {
	ctx := context.TODO()

	cursor, err := store.db.Collection("Deck").Find(ctx, bson.M{"UserID": userID})
	if err != nil {
		store.logger.Fatal(err)
	}
	var decks []entities.Deck
	if err = cursor.All(ctx, &decks); err != nil {
		store.logger.Fatal(err)
	}

	return decks, nil
}

func (store *DeckStore) CreateDeck(deck *entities.Deck) error {
	ctx := context.TODO()

	_, err := store.db.Collection("Deck").InsertOne(ctx, deck)
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}

func (store *DeckStore) UpdateDeck(userID string, deck *entities.Deck) error {
	ctx := context.TODO()

	_, err := store.db.Collection("Deck").
		UpdateOne(ctx, bson.M{"_id": deck.ID, "UserID": userID}, bson.M{"$set": deck})
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}

func (store *DeckStore) DeleteDeck(userID string, id string) error {
	ctx := context.TODO()

	_, err := store.db.Collection("Deck").DeleteOne(ctx, bson.M{"_id": id, "UserID": userID})
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}
