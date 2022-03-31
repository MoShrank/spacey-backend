package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/learning-service/entity"
)

type EventHandler struct {
	logger    logger.LoggerInterface
	usecase   entity.CardEventUsecaseInterface
	validator validator.ValidatorInterface
}

type EventHandlerInterface interface {
	GetLearningCards(c *gin.Context)
	CreateCardEvent(c *gin.Context)
	GetDeckRecallProbabilities(c *gin.Context)
}

func NewEventHandler(
	logger logger.LoggerInterface,
	usecase entity.CardEventUsecaseInterface,
	validator validator.ValidatorInterface,
) EventHandlerInterface {
	return &EventHandler{
		logger:    logger,
		usecase:   usecase,
		validator: validator,
	}
}

func (h *EventHandler) GetLearningCards(c *gin.Context) {
	userID := c.Query("userID")
	ids := c.QueryArray("ids")

	if userID == "" {
		httpconst.WriteBadRequest(c, "userID is required")
		return
	}

	if len(ids) == 0 {
		httpconst.WriteSuccess(c, []entity.CardEventRes{})
		return
	}

	cards, err := h.usecase.GetLearningCards(userID, ids)
	if err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	httpconst.WriteSuccess(c, cards)
}

func (h *EventHandler) CreateCardEvent(c *gin.Context) {
	userID := c.Query("userID")
	if userID == "" {
		httpconst.WriteBadRequest(c, "userID is required")
		return
	}

	var cardEvent entity.CardEventReq
	if err := h.validator.ValidateJSON(c, &cardEvent); err != nil {
		return
	}

	err := h.usecase.CreateCardEvent(userID, &cardEvent)
	if err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	httpconst.WriteCreated(c, nil)
}

func (h *EventHandler) GetDeckRecallProbabilities(c *gin.Context) {
	userID := c.Query("userID")
	if userID == "" {
		httpconst.WriteBadRequest(c, "userID is required")
		return
	}

	var deckData []entity.ProbabilitiesReq
	if err := h.validator.ValidateJSON(c, &deckData); err != nil {
		return
	}

	probabilities, err := h.usecase.CalculateDeckRecallProbabilities(userID, deckData)
	if err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	httpconst.WriteSuccess(c, probabilities)
}
