package auth

import (
	"context"

	"github.com/moshrank/spacey-backend/services/auth/handler"
	"github.com/moshrank/spacey-backend/services/auth/store"
	"github.com/moshrank/spacey-backend/services/auth/usecase"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type AuthService struct {
	router *gin.RouterGroup
}

type AuthServiceInterface interface {
}

func runAddRoutes(
	lifecycle fx.Lifecycle,
	handler handler.HandlerInterface,
	router *gin.RouterGroup,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router.POST("/user", handler.CreateUser)
		router.GET("/user/login", handler.Login)
		router.GET("/user/auth", handler.Authenticate)
		return nil
	}})
}

func NewAuthService(
	router *gin.RouterGroup,
	dbConnection db.DatabaseInterface,
	loggerObj logger.LoggerInterface,
	secretKey string,
) AuthServiceInterface {
	fx.New(
		fx.Provide(func() *gin.RouterGroup { return router }),
		fx.Provide(func() *mongo.Database { return dbConnection.GetDB() }),
		fx.Provide(func() logger.LoggerInterface { return loggerObj }),
		fx.Provide(func() usecase.SecretKey { return usecase.SecretKey{SecretKey: secretKey} }),
		fx.Provide(store.NewStore),
		fx.Provide(usecase.NewUserUseCase),
		fx.Provide(handler.NewHandler),
		fx.Invoke(runAddRoutes),
	).Start(context.TODO())

	return &AuthService{
		router: router,
	}
}
