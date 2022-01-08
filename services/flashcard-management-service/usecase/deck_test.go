package usecase

import (
	"testing"

	"github.com/moshrank/spacey-backend/services/flashcard-management-service/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DeckStoreMock struct {
	mock.Mock
}

func (s *DeckStoreMock) Save(deck *entity.Deck) (string, error) {
	args := s.Called(deck)
	return args.Get(0).(string), args.Error(1)
}

func (s *DeckStoreMock) FindAll(userID string) ([]entity.Deck, error) {
	args := s.Called(userID)
	return args.Get(0).([]entity.Deck), args.Error(1)
}

func (s *DeckStoreMock) FindByID(userID, deckID string) (*entity.Deck, error) {
	args := s.Called(userID, deckID)
	return args.Get(0).(*entity.Deck), args.Error(1)
}

func (s *DeckStoreMock) Update(deck *entity.Deck) error {
	args := s.Called(deck)
	return args.Error(0)
}

func (s *DeckStoreMock) Delete(userID, deckID string) error {
	args := s.Called(userID, deckID)
	return args.Error(0)
}

func TestCreateDeck(t *testing.T) {
	expDeck := entity.DeckRes{
		ID:   "1",
		Name: "Test Deck",
	}

	inpDeck := entity.DeckReq{
		Name: "Test Deck",
	}

	deckStoreMock := new(DeckStoreMock)
	deckStoreMock.On("Save", mock.Anything).Return("1", nil)

	deckUseCase := NewDeckUseCase(deckStoreMock)

	deck, err := deckUseCase.CreateDeck("1", &inpDeck)

	assert.Nil(t, err)
	assert.Equal(t, &expDeck, deck)
}

func TestGetDecks(t *testing.T) {
	expDecks := []entity.DeckRes{
		{
			ID:   "1",
			Name: "Test Deck",
		},
	}

	deckStoreMock := new(DeckStoreMock)
	deckStoreMock.On("FindAll", mock.Anything).Return([]entity.Deck{
		{
			ID:        "1",
			Name:      "Test Deck",
			UserID:    "1",
			CreatedAt: nil,
			UpdatedAt: nil,
			DeletedAt: nil,
		},
	}, nil)

	deckUseCase := NewDeckUseCase(deckStoreMock)

	decks, err := deckUseCase.GetDecks("1")

	assert.Nil(t, err)
	assert.Equal(t, expDecks, decks)
}

func TestGetDeck(t *testing.T) {
	expDeck := entity.DeckRes{

		ID:   "1",
		Name: "Test Deck",
	}

	deckStoreMock := new(DeckStoreMock)
	deckStoreMock.On("FindByID", mock.Anything, mock.Anything).Return(&entity.Deck{

		ID:        "1",
		Name:      "Test Deck",
		UserID:    "1",
		CreatedAt: nil,
		UpdatedAt: nil,
		DeletedAt: nil,
	}, nil)

	deckUseCase := NewDeckUseCase(deckStoreMock)

	decks, err := deckUseCase.GetDeck("1", "1")

	assert.Nil(t, err)
	assert.Equal(t, &expDeck, decks)
}

func TestUpdateDeck(t *testing.T) {
	expDeck := entity.DeckRes{
		ID:   "1",
		Name: "Test Deck",
	}

	inpDeck := entity.DeckReq{
		Name: "Test Deck",
	}

	deckStoreMock := new(DeckStoreMock)
	deckStoreMock.On("Update", mock.Anything).Return(nil)
	deckUseCase := NewDeckUseCase(deckStoreMock)

	deck, err := deckUseCase.UpdateDeck("1", "1", &inpDeck)

	assert.Nil(t, err)
	assert.Equal(t, &expDeck, deck)
}

func TestDeleteDeck(t *testing.T) {
	deckStoreMock := new(DeckStoreMock)
	deckStoreMock.On("Delete", mock.Anything, mock.Anything).Return(nil)
	deckUseCase := NewDeckUseCase(deckStoreMock)

	err := deckUseCase.DeleteDeck("1", "1")

	assert.Nil(t, err)
}
