package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/learning-service/entity"
)

type LearningSessionHandler struct {
	LearningSessionUsecase entity.LearningSessionUsecaseInterface
	logger                 logger.LoggerInterface
	validator              validator.ValidatorInterface
}

type LearningSessionHandlerInterface interface {
	CreateLearningSession(c *gin.Context)
	FinishLearningSession(c *gin.Context)
}

func NewLearningSessionHandler(
	LearningSessionUsecase entity.LearningSessionUsecaseInterface,
	logger logger.LoggerInterface,
	validator validator.ValidatorInterface,
) LearningSessionHandlerInterface {
	return &LearningSessionHandler{
		LearningSessionUsecase: LearningSessionUsecase,
		logger:                 logger,
		validator:              validator,
	}
}

func (h *LearningSessionHandler) CreateLearningSession(c *gin.Context) {
	userID := c.Query("userID")
	if userID == "" {
		httpconst.WriteBadRequest(c, "userID is required")
		return
	}

	var session entity.LearningSessionCreateReq
	if err := h.validator.ValidateJSON(c, &session); err != nil {
		return
	}

	sessionID, err := h.LearningSessionUsecase.CreateLearningSession(userID, &session)
	if err != nil {
		httpconst.WriteInternalServerError(c, err.Error())
		return
	}

	httpconst.WriteCreated(c, &entity.LearningSessionRes{ID: sessionID})
}

func (h *LearningSessionHandler) FinishLearningSession(c *gin.Context) {
	userID := c.Query("userID")

	var session entity.LearningSessionUpdateReq
	if err := h.validator.ValidateJSON(c, &session); err != nil {
		return
	}

	err := h.LearningSessionUsecase.FinishLearningSession(userID, &session)
	if err != nil {
		httpconst.WriteInternalServerError(c, err.Error())
		return
	}

	httpconst.WriteSuccess(c, nil)
}
