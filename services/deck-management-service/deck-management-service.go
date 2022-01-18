package main

import (
	"context"

	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
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
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router := gin.Default()

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

		router.Run(":" + cfg.GetPort())
		return nil
	}})
}

func main() {
	fx.New(
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
