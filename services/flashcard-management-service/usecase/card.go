package usecase

import (
	"time"

	mapper "github.com/PeteProgrammer/go-automapper"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/entity"
)

type CardUseCase struct {
	cardStore entity.CardStoreInterface
}

func NewCardUseCase(cardStore entity.CardStoreInterface) entity.CardUseCaseInterface {
	return &CardUseCase{
		cardStore: cardStore,
	}
}

func (c *CardUseCase) CreateCard(userID string, card *entity.CardReq) (*entity.CardRes, error) {
	var cardDB entity.Card
	mapper.MapLoose(card, &cardDB)

	timestamp := time.Now()
	cardDB.CreatedAt = &timestamp
	cardDB.UpdatedAt = &timestamp
	cardDB.DeletedAt = nil
	cardDB.UserID = userID

	cardID, err := c.cardStore.SaveCard(&cardDB)
	if err != nil {
		return nil, err
	}

	var cardRes entity.CardRes
	mapper.MapLoose(&cardDB, &cardRes)
	cardRes.ID = cardID

	return &cardRes, nil
}

func (c *CardUseCase) GetCardsByDeckID(userID, deckID string) ([]entity.CardRes, error) {
	cards, err := c.cardStore.GetCardsByDeckID(userID, deckID)
	if err != nil {
		return nil, err
	}

	var cardsRes []entity.CardRes
	mapper.Map(cards, &cardsRes)

	return cardsRes, nil
}

func (c *CardUseCase) GetCardByID(userID, cardID string) (*entity.CardRes, error) {
	card, err := c.cardStore.GetCardByID(userID, cardID)
	if err != nil {
		return nil, err
	}

	var cardRes entity.CardRes
	mapper.MapLoose(card, &cardRes)

	return &cardRes, nil
}

//TODO redefine interface to take card ID
//TODO check if all fields or only non nil fields are updated
func (c *CardUseCase) UpdateCard(
	userID, cardID string,
	card *entity.CardReq,
) (*entity.CardRes, error) {
	var cardDB entity.Card
	mapper.MapLoose(card, &cardDB)

	timestamp := time.Now()
	cardDB.UpdatedAt = &timestamp
	cardDB.DeletedAt = nil

	err := c.cardStore.UpdateCard(userID, &cardDB)
	if err != nil {
		return nil, err
	}

	var cardRes entity.CardRes
	mapper.MapLoose(&cardDB, &cardRes)
	cardRes.ID = cardID

	return &cardRes, nil
}

func (c *CardUseCase) DeleteCard(userID, cardID string) error {
<<<<<<< HEAD
	return c.cardStore.DeleteCard(userID, cardID)
=======
	return nil
>>>>>>> 9f50d37... .
}
