package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/deck-management-service/entity"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DeckUseCaseMock struct {
	mock.Mock
}

func (u *DeckUseCaseMock) CreateDeck(userID string, deck *entity.DeckReq) (*entity.DeckRes, error) {
	args := u.Called(userID, deck)
	return args.Get(0).(*entity.DeckRes), args.Error(1)
}

func (u *DeckUseCaseMock) GetDecks(userID string) ([]entity.DeckRes, error) {
	args := u.Called(userID)
	return args.Get(0).([]entity.DeckRes), args.Error(1)
}

func (u *DeckUseCaseMock) GetDeck(userID, DeckID string) (*entity.DeckRes, error) {
	args := u.Called(userID, DeckID)
	return args.Get(0).(*entity.DeckRes), args.Error(1)
}

func (u *DeckUseCaseMock) UpdateDeck(
	userID, DeckID string,
	deck *entity.DeckReq,
) (*entity.DeckRes, error) {
	args := u.Called(userID, DeckID, deck)
	return args.Get(0).(*entity.DeckRes), args.Error(1)
}

func (u *DeckUseCaseMock) DeleteDeck(userID, deckID string) error {
	args := u.Called(userID, deckID)
	return args.Error(0)
}

var validatorObj = validator.NewValidator()

func TestCreateDeck(t *testing.T) {
	tests := []struct {
		testName       string
		body           string
		userID         string
		wantStatusCode int
	}{
		{
			"Valid Deck",
			`{"name": "Test Deck", "description": "Test Description", "color": "#FFFFFF"}`,
			"test_user_id",
			201,
		},
		{
			"Missing Deck Name",
			`{"description": "Test Description", "color": "#FFFFFF"}`,
			"test_user_id",
			400,
		},
		{
			"Missing Color",
			`{"name": "Test Deck", "description": "Test Description"}`,
			"test_user_id",
			400,
		},
		{
			"Missing User ID",
			`{"name": "Test Deck"}`,
			"",
			401,
		},
		{
			"Empty Body",
			"",
			"test_user_id",
			400,
		},
		{
			"Long Deck Name",
			`{"name": "Test Deck with more than 30 characters", "description": "Test Description", "color": "#FFFFFF"}`,
			"test_user_id",
			400,
		},
	}

	deckUseCaseMock := new(DeckUseCaseMock)

	var handler = NewDeckHandler(log.New(), deckUseCaseMock, validatorObj)

	deckUseCaseMock.On("CreateDeck", mock.Anything, mock.Anything).Return(&entity.DeckRes{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"POST",
				"/decks",
				bytes.NewBuffer([]byte(test.body)),
			)
			q := c.Request.URL.Query()
			q.Add("userID", test.userID)
			c.Request.URL.RawQuery = q.Encode()

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

	deckUseCaseMock := new(DeckUseCaseMock)

	var handler = NewDeckHandler(log.New(), deckUseCaseMock, validatorObj)

	deckUseCaseMock.On("GetDecks", mock.Anything).Return([]entity.DeckRes{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				"/decks",
				nil,
			)
			q := c.Request.URL.Query()
			q.Add("userID", test.userID)
			c.Request.URL.RawQuery = q.Encode()

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

	deckUseCaseMock := new(DeckUseCaseMock)

	var handler = NewDeckHandler(log.New(), deckUseCaseMock, validatorObj)

	deckUseCaseMock.On("GetDeck", mock.Anything, mock.Anything).Return(&entity.DeckRes{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				"/decks/:deckID",
				nil,
			)
			c.Params = []gin.Param{
				{
					Key:   "deckID",
					Value: test.deckID,
				},
			}
			q := c.Request.URL.Query()
			q.Add("userID", test.userID)
			c.Request.URL.RawQuery = q.Encode()

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
			`{"name": "Test Deck", "description": "Test Description", "color": "#FFFFFF"}`,
			"test_user_id",
			"test_deck_id",
			200,
		},
		{
			"Invalid Deck",
			`{"name": ""}`,
			"test_user_id",
			"test_deck_id",
			400,
		},
		{
			"Missing User ID",
			`{"name": "Test Deck"}`,
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

	deckUseCaseMock := new(DeckUseCaseMock)

	var handler = NewDeckHandler(log.New(), deckUseCaseMock, validatorObj)

	deckUseCaseMock.On("UpdateDeck", mock.Anything, mock.Anything, mock.Anything).
		Return(&entity.DeckRes{}, nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"PUT",
				"/decks/:deckID",
				bytes.NewBuffer([]byte(test.body)),
			)

			q := c.Request.URL.Query()
			q.Add("userID", test.userID)
			c.Request.URL.RawQuery = q.Encode()

			c.Params = []gin.Param{
				{
					Key:   "deckID",
					Value: test.deckID,
				},
			}

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

	deckUseCaseMock := new(DeckUseCaseMock)

	var handler = NewDeckHandler(log.New(), deckUseCaseMock, validatorObj)

	deckUseCaseMock.On("DeleteDeck", mock.Anything, mock.Anything).Return(nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"DELETE",
				"/decks/:deckID",
				nil,
			)
			c.Params = []gin.Param{
				{
					Key:   "deckID",
					Value: test.deckID,
				},
			}
			q := c.Request.URL.Query()
			q.Add("userID", test.userID)
			c.Request.URL.RawQuery = q.Encode()

			handler.DeleteDeck(c)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}
