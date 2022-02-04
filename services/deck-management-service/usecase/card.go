package usecase

import (
	"time"

	mapper "github.com/PeteProgrammer/go-automapper"
	"github.com/moshrank/spacey-backend/services/deck-management-service/entity"
)

type CardUseCase struct {
	cardStore entity.CardStoreInterface
}

func NewCardUseCase(cardStore entity.CardStoreInterface) entity.CardUseCaseInterface {
	return &CardUseCase{
		cardStore: cardStore,
	}
}

func (c *CardUseCase) CreateCard(
	deckID, userID string,
	card *entity.CardReq,
) (*entity.CardRes, error) {
	var cardDB entity.Card
	mapper.MapLoose(card, &cardDB)

	timestamp := time.Now()
	cardDB.CreatedAt = &timestamp
	cardDB.UpdatedAt = &timestamp
	cardDB.DeletedAt = nil
	cardDB.UserID = userID

	cardID, err := c.cardStore.SaveCard(deckID, userID, &cardDB)
	if err != nil {
		return nil, err
	}

	var cardRes entity.CardRes
	mapper.MapLoose(&cardDB, &cardRes)
	cardRes.ID = cardID

	return &cardRes, nil
}

func (c *CardUseCase) UpdateCard(
	cardID, userID, deckID string,
	card *entity.CardReq,
) (*entity.CardRes, error) {
	var cardDB entity.Card
	mapper.MapLoose(card, &cardDB)

	timestamp := time.Now()
	cardDB.UpdatedAt = &timestamp
	cardDB.DeletedAt = nil

	err := c.cardStore.UpdateCard(cardID, userID, deckID, &cardDB)
	if err != nil {
		return nil, err
	}

	var cardRes entity.CardRes
	mapper.MapLoose(&cardDB, &cardRes)

	return &cardRes, nil
}

func (c *CardUseCase) DeleteCard(userID, deckID, cardID string) error {
	return c.cardStore.DeleteCard(userID, deckID, cardID)
}
