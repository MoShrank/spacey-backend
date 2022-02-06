package usecase

import (
	"time"

	mapper "github.com/PeteProgrammer/go-automapper"
	"github.com/moshrank/spacey-backend/services/deck-management-service/entity"
)

type DeckUseCase struct {
	deckStore entity.DeckStoreInterface
}

func NewDeckUseCase(deckStore entity.DeckStoreInterface) entity.DeckUseCaseInterface {
	return &DeckUseCase{
		deckStore: deckStore,
	}
}

func (u *DeckUseCase) CreateDeck(userID string, deck *entity.DeckReq) (*entity.DeckRes, error) {
	var deckDB entity.Deck
	mapper.MapLoose(deck, &deckDB)

	timestamp := time.Now()
	deckDB.CreatedAt = &timestamp
	deckDB.UpdatedAt = &timestamp
	deckDB.DeletedAt = nil
	deckDB.UserID = userID
	deckDB.Cards = []entity.Card{}

	deckID, err := u.deckStore.Save(&deckDB)
	if err != nil {
		return nil, err
	}

	var deckRes entity.DeckRes
	mapper.MapLoose(&deckDB, &deckRes)
	deckRes.ID = deckID

	return &deckRes, nil
}

func (u *DeckUseCase) GetDecks(userID string) ([]entity.DeckRes, error) {
	decks, err := u.deckStore.FindAll(userID)
	if err != nil {
		return nil, err
	}

	var decksRes []entity.DeckRes
	mapper.MapLoose(decks, &decksRes)

	return decksRes, nil
}

func (u *DeckUseCase) GetDeck(userID, DeckID string) (*entity.DeckRes, error) {
	deck, err := u.deckStore.FindByID(userID, DeckID)
	if err != nil {
		return nil, err
	}

	var deckRes entity.DeckRes
	mapper.MapLoose(deck, &deckRes)

	return &deckRes, nil
}

func (u *DeckUseCase) UpdateDeck(
	userID, DeckID string,
	deck *entity.DeckReq,
) (*entity.DeckRes, error) {
	var deckDB entity.Deck
	mapper.MapLoose(deck, &deckDB)

	timestamp := time.Now()
	deckDB.UpdatedAt = &timestamp

	err := u.deckStore.Update(&deckDB)
	if err != nil {
		return nil, err
	}

	var deckRes entity.DeckRes
	mapper.MapLoose(&deckDB, &deckRes)
	deckRes.ID = DeckID

	return &deckRes, nil
}

func (u *DeckUseCase) DeleteDeck(userID, DeckID string) error {
	return u.deckStore.Delete(userID, DeckID)
}
