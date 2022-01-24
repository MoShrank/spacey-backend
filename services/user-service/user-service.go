package main

import (
	"context"

	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/services/user-service/handler"
	"github.com/moshrank/spacey-backend/services/user-service/store"
	"github.com/moshrank/spacey-backend/services/user-service/usecase"

	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func runServer(
	lifecycle fx.Lifecycle,
	handler handler.HandlerInterface,
	cfg config.ConfigInterface,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router := gin.Default()

		router.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		router.GET("/user", handler.GetUser)
		router.POST("/user", handler.CreateUser)
		router.POST("/login", handler.Login)
		router.GET("/logout", handler.Logout)

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
		fx.Provide(auth.NewJWT),
		fx.Provide(store.NewStore),
		fx.Provide(usecase.NewUserUseCase),
		fx.Provide(handler.NewHandler),
		fx.Invoke(runServer),
	).Start(context.TODO())
}
