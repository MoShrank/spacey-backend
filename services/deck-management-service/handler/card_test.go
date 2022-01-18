package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/deck-management-service/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CardUseCaseMock struct {
	mock.Mock
}

func (m *CardUseCaseMock) CreateCard(userID string, Card *entity.CardReq) (*entity.CardRes, error) {
	args := m.Called(userID, Card)
	return args.Get(0).(*entity.CardRes), args.Error(1)
}

func (m *CardUseCaseMock) GetCardsByDeckID(userID, deckID string) ([]entity.CardRes, error) {
	args := m.Called(userID, deckID)
	return args.Get(0).([]entity.CardRes), args.Error(1)
}

func (m *CardUseCaseMock) GetCardByID(userID, cardID string) (*entity.CardRes, error) {
	args := m.Called(userID, cardID)
	return args.Get(0).(*entity.CardRes), args.Error(1)
}

func (m *CardUseCaseMock) UpdateCard(
	userID, cardID string,
	card *entity.CardReq,
) (*entity.CardRes, error) {
	args := m.Called(userID, card)
	return args.Get(0).(*entity.CardRes), args.Error(1)
}

func (m *CardUseCaseMock) DeleteCard(userID string, cardID string) error {
	args := m.Called(userID, cardID)
	return args.Error(0)
}

func TestCreateCard(t *testing.T) {
	tests := []struct {
		testName       string
		body           string
		userID         string
		wantStatusCode int
	}{
		{
			"Empty Body",
			"",
			"test_user_id",
			400,
		},
		{
			"Valid Card",
			"{\"question\": \"Test Question\", \"answer\": \"Test Answer\", \"deck_id\": \"test_deck_id\"}",
			"test_user_id",
			201,
		},
		{
			"Missing Question",
			"{\"question\": \"\", \"answer\": \"Test Answer\", \"deck_id\": \"test_deck_id\"}",
			"test_user_id",
			400,
		},
		{
			"Missing Answer",
			"{\"question\": \"Test Question\", \"answer\": \"\", \"deck_id\": \"test_deck_id\"}",
			"test_user_id",
			400,
		},
		{
			"Missing Deck ID",
			"{\"question\": \"Test Question\", \"answer\": \"Test Answer\", \"deck_id\": \"\"}",
			"test_user_id",
			400,
		},
		{
			"Missing User ID",
			"{\"question\": \"Test Question\", \"answer\": \"Test Answer\", \"deck_id\": \"test_deck_id\"}",
			"",
			401,
		},
	}

	cardUseCaseMock := new(CardUseCaseMock)

	var handler = NewCardHandler(log.New(), cardUseCaseMock, validator.NewValidator())

	cardUseCaseMock.On("CreateCard", mock.Anything, mock.Anything).Return(&entity.CardRes{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"POST",
				"/flashcards/cards",
				bytes.NewBuffer([]byte(test.body)),
			)

			c.Set("userID", test.userID)

			handler.CreateCard(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestGetCards(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		wantStatusCode int
	}{
		{
			"Valid User ID",
			"test_user_id",
			200,
		},
		{
			"Missing User ID",
			"",
			401,
		},
	}

	cardUseCaseMock := new(CardUseCaseMock)

	var handler = NewCardHandler(log.New(), cardUseCaseMock, validator.NewValidator())

	cardUseCaseMock.On("GetCardsByDeckID", mock.Anything, mock.Anything).
		Return([]entity.CardRes{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				"/flashcards/cards?deck_id=test_deck_id",
				nil,
			)

			c.Set("userID", test.userID)

			handler.GetCards(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestGetCard(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		wantStatusCode int
	}{
		{
			"Valid User ID",
			"test_user_id",
			200,
		},
		{
			"Missing User ID",
			"",
			401,
		},
	}

	cardUseCaseMock := new(CardUseCaseMock)

	var handler = NewCardHandler(log.New(), cardUseCaseMock, validator.NewValidator())

	cardUseCaseMock.On("GetCardByID", mock.Anything, mock.Anything).Return(&entity.CardRes{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// create request with url parameter
			c.Request, _ = http.NewRequest(
				"GET",
				"/flashcards/cards/test_card_id",
				nil,
			)
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: "test",
				},
			}
			c.Set("userID", test.userID)

			handler.GetCard(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestUpdateCard(t *testing.T) {
	tests := []struct {
		testName       string
		body           string
		userID         string
		wantStatusCode int
	}{
		{
			"Empty Body",
			"",
			"test_user_id",
			400,
		},
		{
			"Valid Card",
			"{\"question\": \"Test Question\", \"answer\": \"Test Answer\", \"deck_id\": \"test_deck_id\"}",
			"test_user_id",
			200,
		},
		{
			"Missing Question",
			"{\"question\": \"\", \"answer\": \"Test Answer\", \"deck_id\": \"test_deck_id\"}",
			"test_user_id",
			400,
		},
		{
			"Missing Answer",
			"{\"question\": \"Test Question\", \"answer\": \"\", \"deck_id\": \"test_deck_id\"}",
			"test_user_id",
			400,
		},
		{
			"Missing Deck ID",
			"{\"question\": \"Test Question\", \"answer\": \"Test Answer\", \"deck_id\": \"\"}",
			"test_user_id",
			400,
		},
		{
			"Missing User ID",
			"{\"question\": \"Test Question\", \"answer\": \"Test Answer\", \"deck_id\": \"test_deck_id\"}",
			"",
			401,
		},
	}

	cardUseCaseMock := new(CardUseCaseMock)

	var handler = NewCardHandler(log.New(), cardUseCaseMock, validator.NewValidator())

	cardUseCaseMock.On("UpdateCard", mock.Anything, mock.Anything).Return(&entity.CardRes{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"PUT",
				"/flashcards/cards/test_card_id",
				bytes.NewBuffer([]byte(test.body)),
			)
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: "test",
				},
			}

			c.Set("userID", test.userID)

			handler.UpdateCard(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestDeleteCard(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		wantStatusCode int
	}{
		{
			"Valid User ID",
			"test_user_id",
			200,
		},
		{
			"Missing User ID",
			"",
			401,
		},
	}

	cardStoreMock := new(CardUseCaseMock)

	var handler = NewCardHandler(log.New(), cardStoreMock, validator.NewValidator())

	cardStoreMock.On("DeleteCard", mock.Anything, mock.Anything).Return(nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"DELETE",
				"/flashcards/cards/test_id",
				nil,
			)
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: "test",
				},
			}
			c.Set("userID", test.userID)

			handler.DeleteCard(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}