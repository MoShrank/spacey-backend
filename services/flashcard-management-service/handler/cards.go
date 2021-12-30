package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httperror"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/models"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/store.go"
)

type CardHandlerInterface interface {
	CreateCard(c *gin.Context)
	GetCard(c *gin.Context)
	GetCards(c *gin.Context)
	UpdateCard(c *gin.Context)
	DeleteCard(c *gin.Context)
}

type CardHandler struct {
	logger    logger.LoggerInterface
	cardStore store.CardStoreInterface
	validator validator.ValidatorInterface
}

func NewCardHandler(
	loggerObj logger.LoggerInterface,
	store store.CardStoreInterface,
	validatorObj validator.ValidatorInterface,
) CardHandlerInterface {
	return &CardHandler{
		logger:    loggerObj,
		cardStore: store,
		validator: validatorObj,
	}
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	var card models.Card

	if err := h.validator.ValidateJSON(c, &card); err != nil {
		return
	}

	card.UserID = userID.(string)

	if err := h.cardStore.CreateCard(&card); err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(201, gin.H{
		"message": "Card created",
	})

}

func (h *CardHandler) GetCard(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	cardID := c.Param("id")

	card, err := h.cardStore.GetCard(userID.(string), cardID)
	if err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(200, gin.H{
		"card": card,
	})

}

func (h *CardHandler) GetCards(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	cards, err := h.cardStore.GetCards(userID.(string))
	if err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(200, gin.H{
		"cards": cards,
	})
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	var card models.Card

	if err := h.validator.ValidateJSON(c, &card); err != nil {
		return
	}

	cardID := c.Param("id")

	card.UserID = userID.(string)
	card.ID = cardID

	if err := h.cardStore.UpdateCard(&card); err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(200, gin.H{
		"message": "Card updated",
	})

}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	cardID := c.Param("id")

	if err := h.cardStore.DeleteCard(userID.(string), cardID); err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(200, gin.H{
		"message": "Card deleted",
	})
}
