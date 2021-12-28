package store

import (
	"context"

	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeckStore struct {
	db     *mongo.Database
	logger logger.LoggerInterface
}

type DeckStoreInterface interface {
	GetDeck(userID string, id string) (*models.Deck, error)
	GetDecks(userID string) ([]models.Deck, error)
	CreateDeck(deck *models.Deck) error
	UpdateDeck(userID string, deck *models.Deck) error
	DeleteDeck(userID string, id string) error
}

func NewDeckStore(
	db *mongo.Database,
	loggerObj logger.LoggerInterface,
) DeckStoreInterface {
	return &DeckStore{
		db:     db,
		logger: loggerObj,
	}
}

const DECK_COLLECTION = "decks"

func (store *DeckStore) GetDeck(userID string, id string) (*models.Deck, error) {
	ctx := context.TODO()

	var deck models.Deck
	err := store.db.Collection(DECK_COLLECTION).
		FindOne(ctx, bson.M{"_id": id, "UserID": userID}).
		Decode(&deck)

	if err != nil {
		store.logger.Fatal(err)
	}

	return &deck, nil
}

func (store *DeckStore) GetDecks(userID string) ([]models.Deck, error) {
	ctx := context.TODO()

	cursor, err := store.db.Collection(DECK_COLLECTION).Find(ctx, bson.M{"UserID": userID})
	if err != nil {
		store.logger.Fatal(err)
	}
	var decks []models.Deck
	if err = cursor.All(ctx, &decks); err != nil {
		store.logger.Fatal(err)
	}

	return decks, nil
}

func (store *DeckStore) CreateDeck(deck *models.Deck) error {
	ctx := context.TODO()
	_, err := store.db.Collection(DECK_COLLECTION).InsertOne(ctx, deck)
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}

func (store *DeckStore) UpdateDeck(userID string, deck *models.Deck) error {
	ctx := context.TODO()

	_, err := store.db.Collection(DECK_COLLECTION).
		UpdateOne(ctx, bson.M{"_id": deck.ID, "UserID": userID}, bson.M{"$set": deck})
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}

func (store *DeckStore) DeleteDeck(userID string, id string) error {
	ctx := context.TODO()

	_, err := store.db.Collection(DECK_COLLECTION).
		DeleteOne(ctx, bson.M{"_id": id, "UserID": userID})
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}
