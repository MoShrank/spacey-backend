package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httperror"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/models"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/store.go"
)

type DeckHandler struct {
	logger    logger.LoggerInterface
	deckStore store.DeckStoreInterface
	validator validator.ValidatorInterface
}

type DeckHandlerInterface interface {
	GetDeck(c *gin.Context)
	GetDecks(c *gin.Context)
	CreateDeck(c *gin.Context)
	UpdateDeck(c *gin.Context)
	DeleteDeck(c *gin.Context)
}

func NewDeckHandler(
	loggerObj logger.LoggerInterface,
	deckStore store.DeckStoreInterface,
	validator validator.ValidatorInterface,
) DeckHandlerInterface {
	return &DeckHandler{
		logger:    loggerObj,
		deckStore: deckStore,
		validator: validator,
	}
}

func (h *DeckHandler) GetDeck(c *gin.Context) {
	userID, ok := c.Get("userID")

	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	deckID := c.Param("deckID")

	deck, err := h.deckStore.GetDeck(userID.(string), deckID)
	if err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deck": deck,
	})
}

func (h *DeckHandler) GetDecks(c *gin.Context) {
	userID, ok := c.Get("userID")

	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	decks, err := h.deckStore.GetDecks(userID.(string))
	if err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"decks": decks,
	})

}

func (h *DeckHandler) CreateDeck(c *gin.Context) {
	userID, ok := c.Get("userID")

	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	var deck models.Deck
	if err := h.validator.ValidateJSON(c, &deck); err != nil {
		return
	}

	deck.UserID = userID.(string)

	if err := h.deckStore.CreateDeck(&deck); err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Deck created",
	})

}

func (h *DeckHandler) UpdateDeck(c *gin.Context) {
	userID, ok := c.Get("userID")

	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	var deck models.Deck

	if err := h.validator.ValidateJSON(c, &deck); err != nil {
		return
	}

	deck.UserID = userID.(string)
	deck.ID = c.Param("deckID")

	if err := h.deckStore.UpdateDeck(&deck); err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Deck created",
	})
}

func (h *DeckHandler) DeleteDeck(c *gin.Context) {
	userID, ok := c.Get("userID")

	if !ok || userID == "" {
		httperror.Unauthorized(c)
		return
	}

	deckID := c.Param("deckID")

	if err := h.deckStore.DeleteDeck(userID.(string), deckID); err != nil {
		httperror.DatabaseError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Deck deleted",
	})
}
