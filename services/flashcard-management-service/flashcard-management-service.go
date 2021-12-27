package flashcard

import (
	"context"
	"log"

	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type FlashCardService struct {
	router *gin.RouterGroup
}

type FlashCardServiceInterface interface {
}

func NewFlashCardService(
	router *gin.RouterGroup,
	dbObj db.DatabaseInterface,
	loggerObj logger.LoggerInterface,
) FlashCardServiceInterface {
	ctx := context.TODO()

	app := fx.New(
		fx.Provide(func() db.DatabaseInterface { return dbObj }),
		fx.Provide(func() logger.LoggerInterface { return loggerObj }),
		fx.Provide(handler.NewHandler),
		fx.Invoke(runHttpServer),
	)

	if err := app.Start(ctx); err != nil {
		log.Println(err)
	}

	return &FlashCardService{
		router: router,
	}
}

func runHttpServer(
	lifecycle fx.Lifecycle,
	router *gin.RouterGroup,
	cardHandler handler.CardHandlerInterface,
	deckHandler handler.DeckHandlerInterface,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router.GET("/cards", cardHandler.GetCards)
		router.GET("/cards/:id", cardHandler.GetCard)
		router.POST("/cards", cardHandler.CreateCard)
		router.PUT("/cards/:id", cardHandler.UpdateCard)
		router.DELETE("/cards/:id", cardHandler.DeleteCard)
		router.GET("/decks", deckHandler.GetDecks)
		router.GET("/decks/:id", deckHandler.GetDeck)
		router.POST("/decks", deckHandler.CreateDeck)
		router.PUT("/decks/:id", deckHandler.UpdateDeck)
		router.DELETE("/decks/:id", deckHandler.DeleteDeck)
		return nil
	}})
}
