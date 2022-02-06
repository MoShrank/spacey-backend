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

	card.ID = id.Hex()

	// deckID and userID to objectID
	deckIDObj, err := primitive.ObjectIDFromHex(deckID)
	if err != nil {
		return "", err
	}

	_, err = s.db.UpdateDocument(DECK_COLLECTION, bson.M{
		"_id":     deckIDObj,
		"user_id": userID,
	}, bson.M{"$push": bson.M{"cards": card}})

	return id.Hex(), err
}

func (s *CardStore) UpdateCard(cardID, userID, deckID string, Card *entity.Card) error {
	_, err := s.db.UpdateDocument(DECK_COLLECTION, bson.M{
		"_id":     Card.DeckID,
		"user_id": userID,
		"deck_id": deckID,
		"cards":   bson.M{"$set": bson.M{"$elemMatch": bson.M{"_id": cardID}}},
	}, bson.M{"cards.$": Card})

	return err
}

func (s *CardStore) DeleteCard(userID, deckID, cardID string) error {
	_, err := s.db.UpdateDocument(DECK_COLLECTION, bson.M{
		"_id":     deckID,
		"user_id": userID,
		"cards":   bson.M{"$elemMatch": bson.M{"_id": cardID}},
	}, bson.M{"$pull": "cards.$"})

	return err
}
