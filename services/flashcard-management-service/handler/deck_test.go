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

var validatorObj = validator.NewValidator()

type deckStoreMock struct {
	mock.Mock
}

func (m *deckStoreMock) CreateDeck(deck *models.Deck) error {
	args := m.Called(deck)
	return args.Error(0)
}

func (m *deckStoreMock) GetDecks(userID string) ([]models.Deck, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Deck), args.Error(1)
}

func (m *deckStoreMock) GetDeck(userID string, id string) (*models.Deck, error) {
	args := m.Called(userID, id)
	return args.Get(0).(*models.Deck), args.Error(1)
}

func (m *deckStoreMock) UpdateDeck(deck *models.Deck) error {
	args := m.Called(deck)
	return args.Error(0)
}

func (m *deckStoreMock) DeleteDeck(userID string, id string) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

func TestCreateDeck(t *testing.T) {
	tests := []struct {
		testName       string
		body           string
		userID         string
		wantStatusCode int
	}{
		{
			"Valid Deck",
			"{\"name\": \"Test Deck\"}",
			"test_user_id",
			201,
		},
		{
			"Invalid Deck",
			"{\"name\": \"\"}",
			"test_user_id",
			400,
		},
		{
			"Missing User ID",
			"{\"name\": \"Test Deck\"}",
			"",
			401,
		},
		{
			"Empty Body",
			"",
			"test_user_id",
			400,
		},
	}

	deckStoreMock := new(deckStoreMock)

	var handler = NewDeckHandler(logger.NewLogger(""), deckStoreMock, validatorObj)

	deckStoreMock.On("CreateDeck", mock.Anything).Return(nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"POST",
				"/flashcards/decks",
				bytes.NewBuffer([]byte(test.body)),
			)

			c.Set("userID", test.userID)

			handler.CreateDeck(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestGetDecks(t *testing.T) {
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
			"Invalid User ID",
			"",
			401,
		},
	}

	deckStoreMock := new(deckStoreMock)

	var handler = NewDeckHandler(logger.NewLogger(""), deckStoreMock, validatorObj)

	deckStoreMock.On("GetDecks", mock.Anything).Return([]models.Deck{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				"/flashcards/decks",
				nil,
			)

			c.Set("userID", test.userID)

			handler.GetDecks(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestGetDeck(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		deckID         string
		wantStatusCode int
	}{
		{
			"Valid User ID",
			"test_user_id",
			"test_deck_id",
			200,
		},
		{
			"Invalid User ID",
			"",
			"test_deck_id",
			401,
		},
	}

	deckStoreMock := new(deckStoreMock)

	var handler = NewDeckHandler(logger.NewLogger(""), deckStoreMock, validatorObj)

	deckStoreMock.On("GetDeck", mock.Anything, mock.Anything).Return(&models.Deck{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				"/flashcards/decks/"+test.deckID,
				nil,
			)

			c.Set("userID", test.userID)

			handler.GetDeck(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestUpdateDeck(t *testing.T) {
	tests := []struct {
		testName       string
		body           string
		userID         string
		deckID         string
		wantStatusCode int
	}{
		{
			"Valid Deck",
			"{\"name\": \"Test Deck\"}",
			"test_user_id",
			"test_deck_id",
			201,
		},
		{
			"Invalid Deck",
			"{\"name\": \"\"}",
			"test_user_id",
			"test_deck_id",
			400,
		},
		{
			"Missing User ID",
			"{\"name\": \"Test Deck\"}",
			"",
			"test_deck_id",
			401,
		},
		{
			"Empty Body",
			"",
			"test_user_id",
			"test_deck_id",
			400,
		},
	}

	deckStoreMock := new(deckStoreMock)

	var handler = NewDeckHandler(logger.NewLogger(""), deckStoreMock, validatorObj)

	deckStoreMock.On("UpdateDeck", mock.Anything).Return(nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"PUT",
				"/flashcards/decks/"+test.deckID,
				bytes.NewBuffer([]byte(test.body)),
			)

			c.Set("userID", test.userID)

			handler.UpdateDeck(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestDeleteDeck(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		deckID         string
		wantStatusCode int
	}{
		{
			"Valid Deck",
			"test_user_id",
			"test_deck_id",
			200,
		},
		{
			"Missing User ID",
			"",
			"test_deck_id",
			401,
		},
	}

	deckStoreMock := new(deckStoreMock)

	var handler = NewDeckHandler(logger.NewLogger(""), deckStoreMock, validatorObj)

	deckStoreMock.On("DeleteDeck", mock.Anything, mock.Anything).Return(nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"DELETE",
				"/flashcards/decks/"+test.deckID,
				nil,
			)

			c.Set("userID", test.userID)

			handler.DeleteDeck(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}
