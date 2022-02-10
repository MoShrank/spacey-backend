package main

import (
	"context"

	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/deck-management-service/handler"
	"github.com/moshrank/spacey-backend/services/deck-management-service/store.go"
	"github.com/moshrank/spacey-backend/services/deck-management-service/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func runServer(
	lifecycle fx.Lifecycle,
	cardHandler handler.CardHandlerInterface,
	deckHandler handler.DeckHandlerInterface,
	cfg config.ConfigInterface,
	log logger.LoggerInterface,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router := gin.New()
		router.Use(gin.Logger())
		router.Use(middleware.Recovery())

		router.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		router.GET("decks", deckHandler.GetDecks)
		router.GET("decks/:deckID", deckHandler.GetDeck)
		router.POST("decks", deckHandler.CreateDeck)
		router.PUT("decks/:deckID", deckHandler.UpdateDeck)
		router.DELETE("decks/:deckID", deckHandler.DeleteDeck)

		router.POST("decks/:deckID/cards", cardHandler.CreateCard)
		router.PUT("decks/:deckID/cards/:id", cardHandler.UpdateCard)
		router.DELETE("decks/:deckID/cards/:id", cardHandler.DeleteCard)

		router.Run(":" + cfg.GetPort())
		log.Info("Starting server on port: " + cfg.GetPort())
		return nil
	}})
}

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(config.NewConfig),
		fx.Provide(logger.NewLogger),
		fx.Provide(db.NewDB),
		fx.Provide(validator.NewValidator),
		fx.Provide(store.NewDeckStore),
		fx.Provide(store.NewCardStore),
		fx.Provide(usecase.NewCardUseCase),
		fx.Provide(usecase.NewDeckUseCase),
		fx.Provide(handler.NewCardHandler),
		fx.Provide(handler.NewDeckHandler),
		fx.Invoke(runServer),
	).Start(context.TODO())
}
