package store

import (
	"context"
	"fmt"

	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/learning-service/entity"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventStore struct {
	db     db.DatabaseInterface
	logger logger.LoggerInterface
}

func NewEventStore(
	db db.DatabaseInterface,
	logger logger.LoggerInterface,
) entity.CardEventStoreInterface {
	return &EventStore{
		db:     db,
		logger: logger,
	}
}

const cardEventCollection = "cardEvent"

func (s *EventStore) GetLatestEvents(
	userID string,
	cardIDs []string,
) ([]entity.CardEvent, error) {
	/* mongo db get latest unique events filtered by ids */
	db := s.db.GetDB()
	col := db.Collection(cardEventCollection)

	cur, err := col.Aggregate(context.TODO(), []bson.M{
		{
			"$match": bson.M{
				"userID": userID,
				"cardID": bson.M{
					"$in": cardIDs,
				},
			},
		},
		{
			"$sort": bson.M{
				"createdAt": 1,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{"id": "$cardID"},
				"doc": bson.M{"$last": "$$ROOT"},
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$doc",
			},
		},
	})

	if err != nil {
		err = errors.Wrap(err, "could query latest card events")
		s.logger.Error(err)
		return nil, err
	}

	var events []entity.CardEvent

	err = cur.All(context.TODO(), &events)

	if err != nil {
		err = errors.Wrap(err, "could not decode events")
		s.logger.Error(err)
		return nil, err
	}

	return events, nil
}

func (s *EventStore) GetLatestEvent(userID, cardID string) (*entity.CardEvent, error) {
	options := options.Find()
	options.SetSort(bson.M{"createdAt": -1})
	options.SetLimit(1)

	cur, err := s.db.QueryDocuments(cardEventCollection, bson.M{
		"userID": userID,
		"cardID": cardID,
	}, options)

	if err != nil {
		err = errors.Wrap(err, "could not query latest card event")
		s.logger.Error(err)
		return nil, err
	}

	var event entity.CardEvent
	cur.Next(context.TODO())
	err = cur.Decode(&event)
	if err != nil {
		err = errors.Wrap(err, "could not decode card event")
		s.logger.Error(err)
		return nil, err
	}

	return &event, err
}

func (s *EventStore) CreateCardEvent(event *entity.CardEvent) (string, error) {
	res, err := s.db.CreateDocument(cardEventCollection, event)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("could not create card event: %v", event))
		s.logger.Error(err)
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (s *EventStore) GetCardEventsByDeckIDs(
	userID string,
	deckIDs []string,
) ([]entity.DeckCardEvents, error) {
	/* mongo db get latest unique events filtered by ids */
	db := s.db.GetDB()
	col := db.Collection(cardEventCollection)

	cur, err := col.Aggregate(context.TODO(), []bson.M{
		{
			"$match": bson.M{
				"userID": userID,
				"deckID": bson.M{
					"$in": deckIDs,
				},
			},
		},
		{
			"$sort": bson.M{
				"createdAt": 1,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{"id": "$cardID"},
				"doc": bson.M{"$last": "$$ROOT"},
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$doc",
			},
		},
		{
			"$group": bson.M{
				"_id":        "$deckID",
				"cardEvents": bson.M{"$push": "$$ROOT"},
			},
		},
	})

	if err != nil {
		err = errors.Wrap(err, "could query latest card events")
		s.logger.Error(err)
		return nil, err
	}

	var events []entity.DeckCardEvents

	err = cur.All(context.TODO(), &events)
	if err != nil {
		err = errors.Wrap(err, "could not decode events")
		s.logger.Error(err)
		return nil, err
	}

	return events, nil

}
