package entity

import (
	"time"
)

type LearningSession struct {
	ID         string     `bson:"_id ,omitempty"`
	UserID     string     `bson:"userID"`
	DeckID     string     `bson:"deckID"`
	StartedAt  *time.Time `bson:"startedAt"`
	FinishedAt *time.Time `bson:"finishedAt"`
	Finished   bool       `bson:"finished"`
}

type LearningSessionStoreInterface interface {
	Create(session *LearningSession) (string, error)
	Update(userID, sessionID string, finishedAt *time.Time) error
	GetLearningSessionByDay(userID string, startDate, endDate *time.Time) (*LearningSession, error)
}

type LearningSessionUsecaseInterface interface {
	CreateLearningSession(userID string, session *LearningSessionCreateReq) (string, error)
	FinishLearningSession(userID string, session *LearningSessionUpdateReq) error
}

type LearningSessionCreateReq struct {
	DeckID    string     `json:"deckID"    binding:"required"`
	StartedAt *time.Time `json:"startedAt" binding:"required"`
}

type LearningSessionUpdateReq struct {
	ID         string     `json:"id"         binding:"required"`
	FinishedAt *time.Time `json:"finishedAt" binding:"required"`
}

type LearningSessionRes struct {
	ID string `json:"id"`
}
