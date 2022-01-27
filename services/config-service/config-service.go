package main

import (
	"context"

	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/services/config-service/handler"
	"github.com/moshrank/spacey-backend/services/config-service/store"

	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func runServer(
	lifecycle fx.Lifecycle,
	handler handler.ConfigHandlerInterface,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router := gin.Default()

		router.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		router.GET("/config/:configName", handler.GetConfig)

		router.Run(":" + "80")
		return nil
	}})
}

func main() {
	fx.New(
		fx.Provide(config.NewConfig),
		fx.Provide(logger.NewLogger),
		fx.Provide(db.NewDB),
		fx.Provide(store.NewConfigStore),
		fx.Provide(handler.NewConfigHandler),
		fx.Invoke(runServer),
	).Start(context.TODO())
}
