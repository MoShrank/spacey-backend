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

func NewCardStore(
	db *mongo.Database,
	loggerObj logger.LoggerInterface,
) entity.CardStoreInterface {
	return &CardStore{
		db:     db,
		logger: loggerObj,
	}
}

const CARD_COLLECTION = "cards"

func (store *CardStore) GetCardByID(userID, cardID string) (*entity.Card, error) {
	ctx := context.TODO()

	var Card entity.Card
	err := store.db.Collection(CARD_COLLECTION).
		FindOne(ctx, bson.M{"_id": cardID, "UserID": userID}).
		Decode(&Card)

	if err != nil {
		store.logger.Fatal(err)
	}

	return &Card, nil
}

func (store *CardStore) GetCardsByDeckID(userID, cardID string) ([]entity.Card, error) {
	ctx := context.TODO()

	cursor, err := store.db.Collection(CARD_COLLECTION).
		Find(ctx, bson.M{"UserID": userID, "DeckID": cardID})
	if err != nil {
		store.logger.Fatal(err)
	}
	var Cards []entity.Card
	if err = cursor.All(ctx, &Cards); err != nil {
		store.logger.Fatal(err)
	}

	return Cards, nil
}

func (store *CardStore) SaveCard(Card *entity.Card) (string, error) {
	ctx := context.TODO()
	insertionResult, err := store.db.Collection(CARD_COLLECTION).InsertOne(ctx, Card)
	if err != nil {
		store.logger.Fatal(err)
	}

	return insertionResult.InsertedID.(string), nil
}

func (store *CardStore) UpdateCard(userID string, Card *entity.Card) error {
	ctx := context.TODO()

	_, err := store.db.Collection(CARD_COLLECTION).
		UpdateOne(ctx, bson.M{"_id": Card.ID, "userID": userID}, bson.M{"$set": Card})
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
