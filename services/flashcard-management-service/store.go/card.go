package store

import (
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CARD_COLLECTION = "card"

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

func (s *CardStore) GetCardByID(userID, cardID string) (*entity.Card, error) {
	res := s.db.QueryDocument(
		CARD_COLLECTION,
		map[string]interface{}{"_id": cardID, "UserID": userID},
	)

	var card entity.Card
	err := res.Decode(&card)

	return &card, err
}

func (s *CardStore) GetCardsByDeckID(userID, cardID string) ([]entity.Card, error) {
	res, err := s.db.QueryDocuments(
		CARD_COLLECTION,
		map[string]interface{}{"UserID": userID, "DeckID": cardID},
	)

	if err != nil {
		return nil, err
	}

	var cards []entity.Card
	err = res.Decode(&cards)
	return cards, err
}

func (s *CardStore) SaveCard(Card *entity.Card) (string, error) {
	res, err := s.db.CreateDocument(CARD_COLLECTION, Card)
	if err != nil {
		return "", err
	}

	id, err := res.InsertedID.(primitive.ObjectID).MarshalJSON()

	return string(id[:]), err
}

func (s *CardStore) UpdateCard(userID string, Card *entity.Card) error {
	_, err := s.db.UpdateDocument(CARD_COLLECTION, map[string]interface{}{
		"_id":    Card.ID,
		"UserID": userID,
	}, Card)

	return err
}

func (s *CardStore) DeleteCard(userID string, id string) error {
	_, err := s.db.DeleteDocument(CARD_COLLECTION, map[string]interface{}{
		"_id":    id,
		"UserID": userID,
	})

	return err
}
