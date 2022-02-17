package entity

import "time"

type CardEvent struct {
	ID                         string     `bson:"_id,omitempty"`
	UserID                     string     `bson:"userID"`
	DeckID                     string     `bson:"deckID"`
	CardID                     string     `bson:"cardID"`
	LearningSessionID          string     `bson:"learningSessionID"`
	MemoryHalfLife             float64    `bson:"memoryHalfLife"`
	NumberPracticed            int        `bson:"totalNumberPracticed"`
	NumberCorrect              int        `bson:"totalNumberCorrect"`
	NumberIncorrect            int        `bson:"totalNumberIncorrect"`
	NumberPracticedLastSession int        `bson:"totalNumberPracticedLastSession"`
	NumberCorrectLastSession   int        `bson:"totalNumberCorrectLastSession"`
	NumberIncorrectLastSession int        `bson:"totalNumberIncorrectLastSession"`
	CreatedAt                  *time.Time `bson:"createdAt"`
	StartedAt                  *time.Time `bson:"startedAt"`
	FinishedAt                 *time.Time `bson:"finishedAt"`
}

type CardEventStoreInterface interface {
	GetLatestEvents(userID string, cardIDs []string) ([]CardEvent, error)
	GetLatestEvent(userID, cardID string) (*CardEvent, error)
	CreateCardEvent(event *CardEvent) (string, error)
}

type CardEventUsecaseInterface interface {
	GetLearningCards(userID string, cardIDs []string) ([]CardEventRes, error)
	CreateCardEvent(userID string, event *CardEventReq) error
}

type CardEventReq struct {
	DeckID            string     `json:"deckID"            binding:"required"`
	CardID            string     `json:"cardID"            binding:"required"`
	LearningSessionID string     `json:"learningSessionID" binding:"required"`
	StartedAt         *time.Time `json:"startedAt"         binding:"required"`
	FinishedAt        *time.Time `json:"finishedAt"        binding:"required"`
	Correct           bool       `json:"correct"`
}

type CardEventRes struct {
	CardID            string  `json:"cardID"`
	LearningSessionID string  `json:"learningSessionID"`
	RecallProbability float64 `json:"recallProbability"`
}
