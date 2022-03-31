package main

import (
	"context"

	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/learning-service/handler"
	"github.com/moshrank/spacey-backend/services/learning-service/store"
	"github.com/moshrank/spacey-backend/services/learning-service/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func runServer(
	lifecycle fx.Lifecycle,
	cfg config.ConfigInterface,
	eventHandler handler.EventHandlerInterface,
	sessionHandler handler.LearningSessionHandlerInterface,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router := gin.Default()

		router.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		router.POST("session", sessionHandler.CreateLearningSession)
		router.PUT("session", sessionHandler.FinishLearningSession)
		router.POST("event", eventHandler.CreateCardEvent)
		router.GET("events", eventHandler.GetLearningCards)
		router.POST("probabilities", eventHandler.GetDeckRecallProbabilities)

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
		fx.Provide(store.NewEventStore),
		fx.Provide(store.NewLearningSessionsStore),
		fx.Provide(usecase.NewEventUsecase),
		fx.Provide(usecase.NewLearningSessionUsecase),
		fx.Provide(handler.NewEventHandler),
		fx.Provide(handler.NewLearningSessionHandler),
		fx.Invoke(runServer),
	).Start(context.TODO())
}
