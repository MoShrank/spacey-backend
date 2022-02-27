package store

import (
	"fmt"
	"time"

	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/learning-service/entity"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const learningSessionCollection = "learningSession"

type LearningSessionStore struct {
	db     db.DatabaseInterface
	logger logger.LoggerInterface
}

func NewLearningSessionsStore(
	db db.DatabaseInterface,
	logger logger.LoggerInterface,
) entity.LearningSessionStoreInterface {
	return &LearningSessionStore{
		db:     db,
		logger: logger,
	}
}

func (s *LearningSessionStore) Create(session *entity.LearningSession) (string, error) {
	res, err := s.db.CreateDocument(learningSessionCollection, session)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to create learning session: %v", session))
		s.logger.Error(err)
		return "", err
	}

	id, _ := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), err
}

func (s *LearningSessionStore) Update(userID, sessionID string, finishedAt *time.Time) error {
	objID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		err = errors.Wrap(
			err,
			fmt.Sprintf("failed to convert session id to object id: %v", sessionID),
		)
		s.logger.Error(err)
		return err
	}

	_, err = s.db.UpdateDocument(learningSessionCollection, bson.M{
		"_id":    objID,
		"userID": userID,
	}, bson.M{
		"$set": bson.M{
			"finishedAt": finishedAt,
			"finished":   true,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "failed to update learning session")
		s.logger.Error(err)
		return err
	}

	return nil
}
