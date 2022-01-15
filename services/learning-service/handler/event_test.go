package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/moshrank/spacey-backend/pkg/testingutil"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/learning-service/entity"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EventUsecaseMock struct {
	mock.Mock
}

func (u *EventUsecaseMock) GetLearningCards(
	userID string,
	cardIDs []string,
) ([]entity.CardEventRes, error) {
	args := u.Called(userID, cardIDs)
	return args.Get(0).([]entity.CardEventRes), args.Error(1)
}
func (u *EventUsecaseMock) CreateCardEvent(userID string, event *entity.CardEventReq) error {
	args := u.Called(userID, event)
	return args.Error(0)
}

func TestGetValidLearningCards(t *testing.T) {
	expStatusCode := 200

	u := &EventUsecaseMock{}
	handler := NewEventHandler(log.New(), u, nil)
	u.On("GetLearningCards", mock.Anything, []string{"1", "2"}).Return(
		[]entity.CardEventRes{
			{

				CardID:            "1",
				LearningSessionID: "1",
				RecallProbability: 0.5,
			},
			{
				CardID:            "2",
				LearningSessionID: "2",
				RecallProbability: 0.5,
			},
		}, nil,
	)

	w := httptest.NewRecorder()
	c := testingutil.NewTestingContext(w, "GET", "/events", "").
		AddQueryParameter("userID", "1").
		AddQueryParameter("ids", "1").
		AddQueryParameter("ids", "2")

	handler.GetLearningCards(c.Context)

	assert.Equal(t, expStatusCode, w.Code)
}

func TestCreateCardEvent(t *testing.T) {
	expStatusCode := 201
	body := `{
		"deckID": "1",
		"cardID": "1",
		"learningSessionID": "1",
		"startedAt": "2020-01-01T00:00:00Z",
		"finishedAt": "2020-01-01T00:00:00Z",
		"correct": true
		}`

	u := &EventUsecaseMock{}

	handler := NewEventHandler(log.New(), u, validator.NewValidator())
	u.On("CreateCardEvent", mock.Anything, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	c := testingutil.NewTestingContext(w, "POST", "/events", body).
		AddQueryParameter("userID", "1")

	handler.CreateCardEvent(c.Context)
	assert.Equal(t, expStatusCode, w.Code)
}
