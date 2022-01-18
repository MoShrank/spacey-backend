package usecase

import (
	"testing"

	"github.com/moshrank/spacey-backend/services/deck-management-service/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CardStoreMock struct {
	mock.Mock
}

func (c *CardStoreMock) SaveCard(card *entity.Card) (string, error) {
	args := c.Called(card)
	return args.String(0), args.Error(1)
}

func (c *CardStoreMock) GetCardsByDeckID(userID, deckID string) ([]entity.Card, error) {
	args := c.Called(userID, deckID)
	return args.Get(0).([]entity.Card), args.Error(1)
}

func (c *CardStoreMock) GetCardByID(userID, cardID string) (*entity.Card, error) {
	args := c.Called(userID, cardID)
	return args.Get(0).(*entity.Card), args.Error(1)
}

func (c *CardStoreMock) UpdateCard(userID string, card *entity.Card) error {
	args := c.Called(userID, card)
	return args.Error(0)
}

func (c *CardStoreMock) DeleteCard(userID, cardID string) error {
	args := c.Called(userID, cardID)
	return args.Error(0)
}

func TestCreateCardValidCard(t *testing.T) {
	inpCard := entity.CardReq{
		Question: "Test Question",
		Answer:   "Test Answer",
		DeckID:   "test_deck_id",
	}

	expCard := entity.CardRes{
		ID: "test_card_id",
		CardReq: entity.CardReq{
			Question: "Test Question",
			Answer:   "Test Answer",
			DeckID:   "test_deck_id",
		},
	}

	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("SaveCard", mock.Anything).Return("test_card_id", nil)

	cardUseCase := NewCardUseCase(cardStoreMock)
	card, err := cardUseCase.CreateCard("1", &inpCard)

	assert.Nil(t, err)
	assert.Equal(t, &expCard, card)
}

func TestGetCardsByDeckID(t *testing.T) {
	expCardRes := []entity.CardRes{
		{
			ID: "test_card_id",
			CardReq: entity.CardReq{
				Question: "Test Question",
				Answer:   "Test Answer",
				DeckID:   "test_deck_id",
			},
		},
	}

	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("GetCardsByDeckID", mock.Anything, mock.Anything).Return([]entity.Card{
		{
			ID:        "test_card_id",
			Question:  "Test Question",
			Answer:    "Test Answer",
			DeckID:    "test_deck_id",
			UserID:    "1",
			CreatedAt: nil,
			UpdatedAt: nil,
			DeletedAt: nil,
		},
	}, nil)

	cardUseCase := NewCardUseCase(cardStoreMock)
	card, err := cardUseCase.GetCardsByDeckID("1", "test_deck_id")

	assert.Nil(t, err)
	assert.Equal(t, expCardRes, card)
}

func TestGetCardByID(t *testing.T) {
	expCardRes := entity.CardRes{
		ID: "test_card_id",
		CardReq: entity.CardReq{
			Question: "Test Question",
			Answer:   "Test Answer",
			DeckID:   "test_deck_id",
		},
	}

	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("GetCardByID", mock.Anything, mock.Anything).Return(&entity.Card{
		ID:        "test_card_id",
		Question:  "Test Question",
		Answer:    "Test Answer",
		DeckID:    "test_deck_id",
		UserID:    "1",
		CreatedAt: nil,
		UpdatedAt: nil,
		DeletedAt: nil,
	}, nil)

	cardUseCase := NewCardUseCase(cardStoreMock)
	card, err := cardUseCase.GetCardByID("1", "test_card_id")

	assert.Nil(t, err)
	assert.Equal(t, &expCardRes, card)
}

func TestUpdateCard(t *testing.T) {
	inpCard := entity.CardReq{
		Question: "Test Question",
		Answer:   "Test Answer",
		DeckID:   "test_deck_id",
	}

	expOutCard := entity.CardRes{
		ID: "test_card_id",
		CardReq: entity.CardReq{
			Question: "Test Question",
			Answer:   "Test Answer",
			DeckID:   "test_deck_id",
		},
	}

	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("UpdateCard", mock.Anything, mock.Anything).Return(nil)

	cardUseCase := NewCardUseCase(cardStoreMock)
	newCard, err := cardUseCase.UpdateCard("1", "test_card_id", &inpCard)

	assert.Nil(t, err)
	assert.Equal(t, &expOutCard, newCard)
}

func TestDeleteCard(t *testing.T) {
	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("DeleteCard", mock.Anything, mock.Anything).Return(nil)

	cardUseCase := NewCardUseCase(cardStoreMock)
	err := cardUseCase.DeleteCard("1", "test_card_id")

	assert.Nil(t, err)
}
