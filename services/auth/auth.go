package auth

import (
	"context"

	"github.com/moshrank/spacey-backend/services/auth/handler"
	"github.com/moshrank/spacey-backend/services/auth/store"
	"github.com/moshrank/spacey-backend/services/auth/usecase"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/moshrank/spacey-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type AuthService struct {
	router gin.IRoutes
}

type AuthServiceInterface interface {
}

func runAddRoutes(
	lifecycle fx.Lifecycle,
	handler handler.HandlerInterface,
	router gin.IRoutes,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router.POST("/user", handler.CreateUser)
		router.GET("/user/login", handler.Login)
		router.GET("/user/auth", handler.Authenticate)
		return nil
	}})
}

func NewAuthService(
	router gin.IRoutes,
	dbConnection *mongo.Database,
	loggerObj logger.LoggerInterface,
	secretKey string,
) AuthServiceInterface {
	fx.New(
		fx.Provide(func() gin.IRoutes { return router }),
		fx.Provide(func() *mongo.Database { return dbConnection }),
		fx.Provide(func() logger.LoggerInterface { return loggerObj }),
		fx.Provide(
			func() usecase.SecretKey { return usecase.SecretKey{SecretKey: []byte(secretKey)} },
		),
		fx.Provide(store.NewStore),
		fx.Provide(usecase.NewUserUseCase),
		fx.Provide(handler.NewHandler),
		fx.Invoke(runAddRoutes),
	).Start(context.TODO())

	return &AuthService{
		router: router,
	}
}
