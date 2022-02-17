package usecase

import (
	"time"

	mapper "github.com/PeteProgrammer/go-automapper"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/learning-service/entity"
)

type LearningSessionUsecase struct {
	LearningSessionStore entity.LearningSessionStoreInterface
	logger               logger.LoggerInterface
}

func NewLearningSessionUsecase(
	SessionStore entity.LearningSessionStoreInterface,
	logger logger.LoggerInterface,
) entity.LearningSessionUsecaseInterface {
	return &LearningSessionUsecase{
		LearningSessionStore: SessionStore,
		logger:               logger,
	}
}

func (u *LearningSessionUsecase) CreateLearningSession(
	userID string,
	session *entity.LearningSessionCreateReq,
) (string, error) {
	startDate := time.Date(
		session.StartedAt.Year(),
		session.StartedAt.Month(),
		session.StartedAt.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
	endDate := time.Date(
		session.StartedAt.Year(),
		session.StartedAt.Month(),
		session.StartedAt.Day()+1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	existingSession, _ := u.LearningSessionStore.GetLearningSessionByDay(
		userID,
		&startDate,
		&endDate,
	)

	if existingSession != nil {
		return existingSession.ID, nil
	}

	var learningSessionDBObj entity.LearningSession
	mapper.MapLoose(session, &learningSessionDBObj)

	learningSessionDBObj.UserID = userID
	learningSessionDBObj.Finished = false

	return u.LearningSessionStore.Create(&learningSessionDBObj)
}

func (u *LearningSessionUsecase) FinishLearningSession(
	userID string, session *entity.LearningSessionUpdateReq,
) error {
	err := u.LearningSessionStore.Update(userID, session.ID, session.FinishedAt)
	return err
}
