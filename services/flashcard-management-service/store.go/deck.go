package store

import (
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const DECK_COLLECTION = "decks"

type DeckStore struct {
	db     db.DatabaseInterface
	logger logger.LoggerInterface
}

func NewDeckStore(
	db db.DatabaseInterface,
	loggerObj logger.LoggerInterface,
) entity.DeckStoreInterface {
	return &DeckStore{
		db:     db,
		logger: loggerObj,
	}
}

func (s *DeckStore) FindByID(userID string, id string) (*entity.Deck, error) {
	res := s.db.QueryDocument(DECK_COLLECTION, map[string]interface{}{"_id": id, "UserID": userID})
	var deck entity.Deck
	err := res.Decode(&deck)

	return &deck, err
}

func (s *DeckStore) FindAll(userID string) ([]entity.Deck, error) {
	res, err := s.db.QueryDocuments(DECK_COLLECTION, map[string]interface{}{"UserID": userID})
	if err != nil {
		return nil, err
	}

	var decks []entity.Deck
	err = res.Decode(&decks)
	return decks, err

}

func (s *DeckStore) Save(deck *entity.Deck) (string, error) {
	res, err := s.db.CreateDocument(DECK_COLLECTION, deck)
	if err != nil {
		return "", err
	}

	id, err := res.InsertedID.(primitive.ObjectID).MarshalJSON()

	return string(id[:]), err
}

func (s *DeckStore) Update(deck *entity.Deck) error {
	_, err := s.db.UpdateDocument(DECK_COLLECTION, map[string]interface{}{}, deck)

	return err
}

func (s *DeckStore) Delete(userID string, id string) error {
	_, err := s.db.DeleteDocument(
		DECK_COLLECTION,
		map[string]interface{}{"_id": id, "UserID": userID},
	)

	return err
}
