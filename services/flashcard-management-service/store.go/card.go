package store

import (
	"context"

	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardStore struct {
	db     *mongo.Database
	logger logger.LoggerInterface
}

type CardStoreInterface interface {
	GetCard(userID string, id string) (*entity.Card, error)
	GetCards(userID string) ([]entity.Card, error)
	CreateCard(Card *entity.Card) error
	UpdateCard(Card *entity.Card) error
	DeleteCard(userID string, id string) error
}

func NewCardStore(
	db *mongo.Database,
	loggerObj logger.LoggerInterface,
) CardStoreInterface {
	return &CardStore{
		db:     db,
		logger: loggerObj,
	}
}

const CARD_COLLECTION = "cards"

func (store *CardStore) GetCard(userID string, id string) (*entity.Card, error) {
	ctx := context.TODO()

	var Card entity.Card
	err := store.db.Collection(CARD_COLLECTION).
		FindOne(ctx, bson.M{"_id": id, "UserID": userID}).
		Decode(&Card)

	if err != nil {
		store.logger.Fatal(err)
	}

	return &Card, nil
}

func (store *CardStore) GetCards(userID string) ([]entity.Card, error) {
	ctx := context.TODO()

	cursor, err := store.db.Collection(CARD_COLLECTION).Find(ctx, bson.M{"UserID": userID})
	if err != nil {
		store.logger.Fatal(err)
	}
	var Cards []entity.Card
	if err = cursor.All(ctx, &Cards); err != nil {
		store.logger.Fatal(err)
	}

	return Cards, nil
}

func (store *CardStore) CreateCard(Card *entity.Card) error {
	ctx := context.TODO()
	_, err := store.db.Collection(CARD_COLLECTION).InsertOne(ctx, Card)
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}

func (store *CardStore) UpdateCard(Card *entity.Card) error {
	ctx := context.TODO()

	_, err := store.db.Collection(CARD_COLLECTION).
		UpdateOne(ctx, bson.M{"_id": Card.ID}, bson.M{"$set": Card})
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}

func (store *CardStore) DeleteCard(userID string, id string) error {
	ctx := context.TODO()

	_, err := store.db.Collection(CARD_COLLECTION).
		DeleteOne(ctx, bson.M{"_id": id, "UserID": userID})
	if err != nil {
		store.logger.Fatal(err)
	}

	return nil
}
