package handler

import "github.com/gin-gonic/gin"

type CardHandlerInterface interface {
	CreateCard(ctx *gin.Context)
	GetCard(ctx *gin.Context)
	GetCards(ctx *gin.Context)
	UpdateCard(ctx *gin.Context)
	DeleteCard(ctx *gin.Context)
}

type Card struct {
}

func NewHandler() CardHandlerInterface {
	return &Card{}
}

func (c *Card) CreateCard(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Create Card",
	})
}

func (c *Card) GetCard(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Get Card",
	})
}

func (c *Card) GetCards(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Get Cards",
	})
}

func (c *Card) UpdateCard(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Update Card",
	})
}

func (c *Card) DeleteCard(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Delete Card",
	})
}
