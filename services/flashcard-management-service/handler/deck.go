package handler

import "github.com/gin-gonic/gin"

type DeckHandler struct {
}

type DeckHandlerInterface interface {
	GetDeck(c *gin.Context)
	GetDecks(c *gin.Context)
	CreateDeck(c *gin.Context)
	UpdateDeck(c *gin.Context)
	DeleteDeck(c *gin.Context)
}

func NewDeckHandler() DeckHandlerInterface {
	return &DeckHandler{}
}

func (d *DeckHandler) GetDeck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get Deck",
	})
}

func (d *DeckHandler) GetDecks(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get Decks",
	})
}

func (d *DeckHandler) CreateDeck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Create Deck",
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
