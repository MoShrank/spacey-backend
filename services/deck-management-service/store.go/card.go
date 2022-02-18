package store

import (
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/deck-management-service/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CardStore struct {
	db     db.DatabaseInterface
	logger logger.LoggerInterface
}

func NewCardStore(
	db db.DatabaseInterface,
	loggerObj logger.LoggerInterface,
) entity.CardStoreInterface {
	return &CardStore{
		db:     db,
		logger: loggerObj,
	}
}

func (s *CardStore) SaveCard(deckID, userID string, card *entity.Card) (string, error) {
	id := primitive.NewObjectID()

	deckIDObj, err := primitive.ObjectIDFromHex(deckID)
	if err != nil {
		return "", err
	}

	_, err = s.db.UpdateDocument(DECK_COLLECTION, bson.M{
		"_id":     deckIDObj,
		"user_id": userID,
	}, bson.M{"$push": bson.M{"cards": bson.M{
		"_id":        id,
		"question":   card.Question,
		"answer":     card.Answer,
		"user_id":    userID,
		"deck_id":    deckID,
		"created_at": card.CreatedAt,
		"updated_at": card.UpdatedAt,
		"deleted_at": card.DeletedAt,
	}}})

	return id.Hex(), err
}

func (s *CardStore) UpdateCard(cardID, userID, deckID string, card *entity.Card) error {
	deckObjID, err := primitive.ObjectIDFromHex(deckID)
	if err != nil {
		return err
	}

	cardObjID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		return err
	}

	_, err = s.db.UpdateDocument(DECK_COLLECTION, bson.M{
		"_id":       deckObjID,
		"user_id":   userID,
		"cards._id": cardObjID,
	}, bson.M{"$set": bson.M{"cards.$.question": card.Question, "cards.$.answer": card.Answer, "cards.$.updated_at": card.UpdatedAt}})

	return err
}

func (s *CardStore) DeleteCard(userID, deckID, cardID string) error {
	deckObjID, err := primitive.ObjectIDFromHex(deckID)
	if err != nil {
		return err
	}

	cardObjID, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		return err
	}

	_, err = s.db.UpdateDocument(DECK_COLLECTION, bson.M{
		"_id":     deckObjID,
		"user_id": userID,
	}, bson.M{"$pull": bson.M{"cards": bson.M{"_id": cardObjID}}})

	return err
}
