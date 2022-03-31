package usecase

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/moshrank/spacey-backend/services/learning-service/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EventStoreMock struct {
	mock.Mock
}

func (s *EventStoreMock) GetLatestEvents(
	userID string,
	cardIDs []string,
) ([]entity.CardEvent, error) {
	args := s.Called(userID, cardIDs)
	return args.Get(0).([]entity.CardEvent), args.Error(1)
}

func (s *EventStoreMock) GetLatestEvent(
	userID string,
	cardIDs string,
) (*entity.CardEvent, error) {
	args := s.Called(userID, cardIDs)
	return args.Get(0).(*entity.CardEvent), args.Error(1)
}

func (s *EventStoreMock) CreateCardEvent(event *entity.CardEvent) (string, error) {
	args := s.Called(event)
	return args.String(0), args.Error(1)
}

func (s *EventStoreMock) GetCardEventsByDeckIDs(
	userID string,
	deckIDs []string,
) ([]entity.DeckCardEvents, error) {
	args := s.Called(userID, deckIDs)
	return args.Get(0).([]entity.DeckCardEvents), args.Error(1)
}

func TestCalculateRecallProbability(t *testing.T) {
	tests := []struct {
		testName string
		timelag  float64
		h        float64
		expProp  float64
	}{
		{
			"Valid",
			1,
			1,
			0.5,
		},
		{
			"Invalid Time Lag",
			-1,
			1,
			0,
		},
		{
			"Invalid H",
			1,
			0,
			0,
		},
	}

	eventStoreMock := new(EventStoreMock)
	usecase := EventUsecase{
		store:  eventStoreMock,
		logger: log.New(),
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			prop := usecase.calculateRecallProbability(test.timelag, test.h)
			assert.Equal(t, test.expProp, prop)
		})
	}
}
