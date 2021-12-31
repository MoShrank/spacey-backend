package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CardStoreMock struct {
	mock.Mock
}

func (m *CardStoreMock) CreateCard(Card *models.Card) error {
	args := m.Called(Card)
	return args.Error(0)
}

func (m *CardStoreMock) GetCards(userID string) ([]models.Card, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Card), args.Error(1)
}

func (m *CardStoreMock) GetCard(userID string, id string) (*models.Card, error) {
	args := m.Called(userID, id)
	return args.Get(0).(*models.Card), args.Error(1)
}

func (m *CardStoreMock) UpdateCard(Card *models.Card) error {
	args := m.Called(Card)
	return args.Error(0)
}

func (m *CardStoreMock) DeleteCard(userID string, id string) error {
	args := m.Called(userID, id)
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

	cardStoreMock := new(CardStoreMock)

	var handler = NewCardHandler(logger.NewLogger(""), cardStoreMock, validator.NewValidator())

	cardStoreMock.On("CreateCard", mock.Anything).Return(nil)

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

	cardStoreMock := new(CardStoreMock)

	var handler = NewCardHandler(logger.NewLogger(""), cardStoreMock, validator.NewValidator())

	cardStoreMock.On("GetCard", mock.Anything, mock.Anything).Return(&models.Card{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				"/flashcards/cards/test_id",
				nil,
			)

			c.Set("userID", test.userID)

			handler.GetCard(c)

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

	cardStoreMock := new(CardStoreMock)

	var handler = NewCardHandler(logger.NewLogger(""), cardStoreMock, validator.NewValidator())

	cardStoreMock.On("GetCards", mock.Anything, mock.Anything).Return([]models.Card{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				"/flashcards/cards",
				nil,
			)

			c.Set("userID", test.userID)

			handler.GetCards(c)

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

	cardStoreMock := new(CardStoreMock)

	var handler = NewCardHandler(logger.NewLogger(""), cardStoreMock, validator.NewValidator())

	cardStoreMock.On("UpdateCard", mock.Anything).Return(nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"PUT",
				"/flashcards/cards/test_id",
				bytes.NewBuffer([]byte(test.body)),
			)

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

	cardStoreMock := new(CardStoreMock)

	var handler = NewCardHandler(logger.NewLogger(""), cardStoreMock, validator.NewValidator())

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

			c.Set("userID", test.userID)

			handler.DeleteCard(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}
