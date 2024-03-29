package main

import (
	"context"

	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/services/user-service/external"
	"github.com/moshrank/spacey-backend/services/user-service/handler"
	"github.com/moshrank/spacey-backend/services/user-service/store"
	"github.com/moshrank/spacey-backend/services/user-service/usecase"

	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/pkg/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type UserServiceConfig struct {
	PORT                string `json:"port"`
	LOGLEVEL            string `json:"loglevel"`
	GRAYLOG_CONNECTION  string `json:"graylog_connection"`
	DB_NAME             string `json:"db_name"`
	MONGO_DB_CONNECTION string `json:"mongo_db_connection"`
	JWT_SECRET          string `json:"jwt_secret"`
	MAX_AGE_AUTH        int    `json:"max_age_auth"`
}

func runServer(
	lifecycle fx.Lifecycle,
	handler handler.HandlerInterface,
	cfg config.ConfigInterface,
	log logger.LoggerInterface,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router := gin.New()
		router.Use(middleware.Logger(log))
		router.Use(middleware.Recovery())

		router.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		router.GET("/user", handler.GetUser)
		router.POST("/user", handler.CreateUser)
		router.POST("/login", handler.Login)
		router.GET("/logout", handler.Logout)
		router.GET("/validate", handler.SendValidationEmail)
		router.POST("/validate", handler.Validate)

		log.Info("Starting server on port: " + cfg.GetPort())
		router.Run(":" + cfg.GetPort())
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
		fx.Provide(auth.NewJWT),
		fx.Provide(external.NewEmailSender),
		fx.Provide(store.NewStore),
		fx.Provide(usecase.NewUserUseCase),
		fx.Provide(handler.NewHandler),
		fx.Invoke(runServer),
	).Start(context.TODO())
}
