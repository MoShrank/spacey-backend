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

func (c *CardStoreMock) SaveCard(deckID, userID string, card *entity.Card) (string, error) {
	args := c.Called(deckID, userID, card)
	return args.String(0), args.Error(1)
}

func (c *CardStoreMock) UpdateCard(cardID, userID, deckID string, card *entity.Card) error {
	args := c.Called(cardID, userID, deckID, card)
	return args.Error(0)
}

func (c *CardStoreMock) DeleteCard(userID, deckID, cardID string) error {
	args := c.Called(userID, deckID, cardID)
	return args.Error(0)
}

func (c *CardStoreMock) SaveCards(
	deckID, userID string,
	cards []entity.Card,
) ([]string, error) {
	args := c.Called(deckID, userID, cards)
	return args.Get(0).([]string), args.Error(1)
}

func TestCreateCardValidCard(t *testing.T) {
	inpCard := entity.CardReq{
		Question: "Test Question",
		Answer:   "Test Answer",
		DeckID:   "test_deck_id",
	}

	expCard := entity.CardRes{
		ID:       "test_card_id",
		Question: "Test Question",
		Answer:   "Test Answer",
		DeckID:   "test_deck_id",
	}

	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("SaveCard", mock.Anything, mock.Anything, mock.Anything).
		Return("test_card_id", nil)

	cardUseCase := NewCardUseCase(cardStoreMock)
	card, err := cardUseCase.CreateCard("1", "1", &inpCard)

	assert.Nil(t, err)
	assert.Equal(t, &expCard, card)
}

func TestUpdateCard(t *testing.T) {
	inpCard := entity.CardReq{
		ID:       "test_card_id",
		Question: "Test Question",
		Answer:   "Test Answer",
		DeckID:   "test_deck_id",
	}

	expOutCard := entity.CardRes{
		ID:       "test_card_id",
		Question: "Test Question",
		Answer:   "Test Answer",
		DeckID:   "test_deck_id",
	}

	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("UpdateCard", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	cardUseCase := NewCardUseCase(cardStoreMock)
	newCard, err := cardUseCase.UpdateCard("1", "1", "test_card_id", &inpCard)

	assert.Nil(t, err)
	assert.Equal(t, &expOutCard, newCard)
}

func TestDeleteCard(t *testing.T) {
	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("DeleteCard", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	cardUseCase := NewCardUseCase(cardStoreMock)
	err := cardUseCase.DeleteCard("1", "1", "test_card_id")

	assert.Nil(t, err)
}

func TestCreateCards(t *testing.T) {
	inpCards := []entity.CardReq{
		{
			Question: "Test Question",
			Answer:   "Test Answer",
			DeckID:   "test_deck_id",
		},
	}

	expOutputCards := []entity.CardRes{
		{
			ID:       "test_card_id",
			Question: "Test Question",
			Answer:   "Test Answer",
			DeckID:   "test_deck_id",
		},
	}

	cardStoreMock := new(CardStoreMock)
	cardStoreMock.On("SaveCards", mock.Anything, mock.Anything, mock.Anything).
		Return([]string{"test_card_id"}, nil)

	cardUseCase := NewCardUseCase(cardStoreMock)

	cards, err := cardUseCase.CreateCards("test_deck_id", "1", inpCards)

	assert.Nil(t, err)
	assert.Equal(t, expOutputCards, cards)
}
