package usecase

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/learning-service/entity"
)

type EventUsecase struct {
	logger logger.LoggerInterface
	store  entity.CardEventStoreInterface
}

func NewEventUsecase(
	logger logger.LoggerInterface,
	store entity.CardEventStoreInterface,
) entity.CardEventUsecaseInterface {
	return &EventUsecase{
		logger: logger,
		store:  store,
	}
}

func (u *EventUsecase) calculateRecallProbability(timeLag float64, h float64) float64 {
	if h <= 0 {
		u.logger.Warn(fmt.Sprintf("cannot calculate recall probability with h <= 0. h: %v", h))
		return 0
	}

	if timeLag < 0 {

		u.logger.Error(
			fmt.Sprintf(
				"cannot calculate recall probability with timeLag <= 0. timeLag: %v",
				timeLag,
			),
		)
		return 0
	}

	return math.Pow(2, (-timeLag / h))
}

func (u *EventUsecase) getTimelag(createdAt *time.Time) int {
	now := time.Now()
	timelag := now.Sub(*createdAt).Hours() / 24
	return int(timelag)
}

func (u *EventUsecase) remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (u *EventUsecase) removeIDFromArray(ids []string, targetID string) []string {
	newIDs := ids
	for idx, id := range ids {
		if targetID == id {
			newIDs = u.remove(ids, idx)
		}
	}

	return newIDs
}

func (u *EventUsecase) appendEmptyEvents(
	ids []string,
	events []entity.CardEventRes,
) []entity.CardEventRes {
	newEvents := events

	for _, id := range ids {
		newEvents = append(
			newEvents,
			entity.CardEventRes{
				CardID:            id,
				RecallProbability: 0,
			},
		)
	}

	return newEvents
}

func (u *EventUsecase) sortByRecallProbability(events []entity.CardEventRes) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].RecallProbability < events[j].RecallProbability
	})
}

func (u *EventUsecase) GetLearningCards(
	userID string,
	ids []string,
) ([]entity.CardEventRes, error) {
	threshold := 0.5

	res, err := u.store.GetLatestEvents(userID, ids)
	if err != nil {
		return nil, err
	}

	var cardEventRes []entity.CardEventRes
	fmt.Println(res)
	for _, e := range res {
		timelag := u.getTimelag(e.CreatedAt)
		fmt.Println(timelag)
		fmt.Println(e.MemoryHalfLife)
		recallProbability := u.calculateRecallProbability(
			float64(timelag),
			float64(e.MemoryHalfLife),
		)

		fmt.Println(recallProbability)

		if recallProbability <= threshold {
			cardEventRes = append(
				cardEventRes,
				entity.CardEventRes{
					CardID:            e.CardID,
					RecallProbability: recallProbability,
				},
			)
		}

		ids = u.removeIDFromArray(ids, e.CardID)
	}

	cardEventRes = u.appendEmptyEvents(ids, cardEventRes)
	u.sortByRecallProbability(cardEventRes)

	return cardEventRes, nil
}

func (u *EventUsecase) calculateHalfLife(currentHalflife float64, correct bool) float64 {
	newHalflife := 0.

	if correct {
		if currentHalflife < 1. {
			newHalflife = 1.
		} else {
			newHalflife = currentHalflife * 2.
		}
	} else {

		if currentHalflife > 1. {
			newHalflife = currentHalflife / 2.
		} else {
			newHalflife = 0.
		}
	}

	return newHalflife
}

func (u *EventUsecase) createEmptyCardEvent(
	cardID, userID, deckID, learningSessionID string,
	startedAt, finishedAt *time.Time,
	correct bool,
) *entity.CardEvent {
	now := time.Now()

	newHalflife := u.calculateHalfLife(0., correct)

	numberCorrect := 0
	numberIncorrect := 0
	numberCorrectLastSession := 0
	numberIncorrectLastSession := 0

	if correct {
		numberCorrect++
		numberCorrectLastSession++
	} else {
		numberIncorrect++
		numberIncorrectLastSession++
	}

	cardEvent := entity.CardEvent{
		CardID:                     cardID,
		UserID:                     userID,
		DeckID:                     deckID,
		LearningSessionID:          learningSessionID,
		MemoryHalfLife:             newHalflife,
		NumberPracticed:            1,
		NumberCorrect:              numberCorrect,
		NumberIncorrect:            numberIncorrect,
		NumberPracticedLastSession: 1,
		NumberCorrectLastSession:   numberCorrectLastSession,
		NumberIncorrectLastSession: numberIncorrectLastSession,
		CreatedAt:                  &now,
		StartedAt:                  startedAt,
		FinishedAt:                 finishedAt,
	}

	return &cardEvent
}

func (u *EventUsecase) CreateCardEvent(
	userID string,
	cardEventReq *entity.CardEventReq,
) error {
	var newCardEvent entity.CardEvent

	cardEvent, err := u.store.GetLatestEvent(userID, cardEventReq.CardID)
	if err != nil {
		newCardEvent = *u.createEmptyCardEvent(
			cardEventReq.CardID,
			userID,
			cardEventReq.DeckID,
			cardEventReq.LearningSessionID,
			cardEventReq.StartedAt,
			cardEventReq.FinishedAt,
			cardEventReq.Correct,
		)
	} else {

		newHalflife := u.calculateHalfLife(cardEvent.MemoryHalfLife, cardEventReq.Correct)

		numberCorrect := cardEvent.NumberCorrect
		numberIncorrect := cardEvent.NumberIncorrect
		numberCorrectLastSession := cardEvent.NumberCorrectLastSession
		numberIncorrectLastSession := cardEvent.NumberIncorrectLastSession

		if cardEventReq.Correct {
			numberCorrect++
			numberCorrectLastSession++
		} else {
			numberIncorrect++
			numberIncorrectLastSession++
		}

		now := time.Now()

		newCardEvent = entity.CardEvent{
			CardID:                     cardEvent.CardID,
			UserID:                     cardEvent.UserID,
			DeckID:                     cardEvent.DeckID,
			LearningSessionID:          cardEvent.LearningSessionID,
			MemoryHalfLife:             newHalflife,
			NumberPracticed:            cardEvent.NumberPracticed + 1,
			NumberCorrect:              numberCorrect,
			NumberIncorrect:            numberIncorrect,
			NumberPracticedLastSession: cardEvent.NumberPracticedLastSession + 1,
			NumberCorrectLastSession:   numberCorrectLastSession,
			NumberIncorrectLastSession: numberIncorrectLastSession,
			CreatedAt:                  &now,
			StartedAt:                  cardEvent.StartedAt,
			FinishedAt:                 cardEvent.FinishedAt,
		}
	}

	_, err = u.store.CreateCardEvent(&newCardEvent)

	return err
}
