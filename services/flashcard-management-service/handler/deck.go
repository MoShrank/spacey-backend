package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/models"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/store.go"
)

type DeckHandler struct {
	logger    logger.LoggerInterface
	deckStore store.DeckStoreInterface
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
) DeckHandlerInterface {
	return &DeckHandler{
		logger:    loggerObj,
		deckStore: deckStore,
	}
}

func (d *DeckHandler) GetDeck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get Deck",
	})
}

func (d *DeckHandler) GetDecks(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		d.logger.Error("userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing userID",
		})
		return
	}

	decks, err := d.deckStore.GetDecks(userID.(string))
	if err != nil {
		d.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting decks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"decks": decks,
	})

}

func (h *DeckHandler) CreateDeck(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		h.logger.Error("userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing userID",
		})
		return
	}

	var deck models.Deck
	err := c.BindJSON(&deck)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error binding deck",
		})
		return
	}

	deck.UserID = userID.(string)

	h.logger.Info(userID)

	h.logger.Info(deck.UserID)

	if err := h.deckStore.CreateDeck(&deck); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating deck",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Deck created",
	})

}

func (d *DeckHandler) UpdateDeck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Update Deck",
	})
}

func (d *DeckHandler) DeleteDeck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Delete Deck",
	})
}
