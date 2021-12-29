package user

import (
	"context"

	"github.com/moshrank/spacey-backend/services/user-service/handler"
	"github.com/moshrank/spacey-backend/services/user-service/store"
	"github.com/moshrank/spacey-backend/services/user-service/usecase"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/moshrank/spacey-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type UserService struct {
	router gin.IRoutes
}

type UserServiceInterface interface {
}

func runAddRoutes(
	lifecycle fx.Lifecycle,
	handler handler.HandlerInterface,
	router gin.IRoutes,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router.POST("/user", handler.CreateUser)
		router.GET("/user/login", handler.Login)
		return nil
	}})
}

func NewUserService(
	router gin.IRoutes,
	dbConnection *mongo.Database,
	loggerObj logger.LoggerInterface,
	secretKey string,
) UserServiceInterface {
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

	return &UserService{
		router: router,
	}
}
