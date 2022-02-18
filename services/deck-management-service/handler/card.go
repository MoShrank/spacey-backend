package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/deck-management-service/entity"
)

type CardHandler struct {
	logger      logger.LoggerInterface
	cardUseCase entity.CardUseCaseInterface
	validator   validator.ValidatorInterface
}

type CardHandlerInterface interface {
	CreateCard(c *gin.Context)
	UpdateCard(c *gin.Context)
	DeleteCard(c *gin.Context)
}

func NewCardHandler(
	loggerObj logger.LoggerInterface,
	cardUseCase entity.CardUseCaseInterface,
	validatorObj validator.ValidatorInterface,
) CardHandlerInterface {
	return &CardHandler{
		logger:      loggerObj,
		cardUseCase: cardUseCase,
		validator:   validatorObj,
	}
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	deckID := c.Param("deckID")

	var card entity.CardReq
	if err := h.validator.ValidateJSON(c, &card); err != nil {
		return
	}

	cardRes, err := h.cardUseCase.CreateCard(deckID, userID, &card)
	if err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	httpconst.WriteCreated(c, cardRes)
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	var card entity.CardReq
	if err := h.validator.ValidateJSON(c, &card); err != nil {
		return
	}

	deckID := c.Param("deckID")
	cardID := c.Param("id")

	cardRes, err := h.cardUseCase.UpdateCard(cardID, userID, deckID, &card)
	if err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	httpconst.WriteSuccess(c, cardRes)
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	deckID := c.Param("deckID")
	cardID := c.Param("id")

	if err := h.cardUseCase.DeleteCard(userID, deckID, cardID); err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	httpconst.WriteSuccess(c, "Deleted Card")
}
