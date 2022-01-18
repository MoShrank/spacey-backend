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
	GetCard(c *gin.Context)
	GetCards(c *gin.Context)
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
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	var card entity.CardReq
	if err := h.validator.ValidateJSON(c, &card); err != nil {
		return
	}

	cardRes, err := h.cardUseCase.CreateCard(userID.(string), &card)
	if err != nil {
		httpconst.WriteInternalServerError(c)
		return
	}

	httpconst.WriteCreated(c, cardRes)
}

func (h *CardHandler) GetCards(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	deckID := c.Query("deck_id")
	if deckID == "" {
		httpconst.WriteBadRequest(c)
		return
	}

	cards, err := h.cardUseCase.GetCardsByDeckID(userID.(string), deckID)
	if err != nil {
		httpconst.WriteInternalServerError(c)
		return
	}

	httpconst.WriteSuccess(c, cards)
}

func (h *CardHandler) GetCard(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	cardID := c.Param("id")
	if cardID == "" {
		httpconst.WriteBadRequest(c)
		return
	}

	card, err := h.cardUseCase.GetCardByID(userID.(string), cardID)
	if err != nil {
		httpconst.WriteInternalServerError(c)
		return
	}

	httpconst.WriteSuccess(c, card)
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	var card entity.CardReq
	if err := h.validator.ValidateJSON(c, &card); err != nil {
		return
	}

	cardID := c.Param("id")
	if cardID == "" {
		httpconst.WriteBadRequest(c)
		return
	}

	cardRes, err := h.cardUseCase.UpdateCard(userID.(string), cardID, &card)
	if err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, cardRes)
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	cardID := c.Param("id")
	if cardID == "" {
		httpconst.WriteBadRequest(c)
		return
	}

	if err := h.cardUseCase.DeleteCard(userID.(string), cardID); err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, "Deleted Card")
}
