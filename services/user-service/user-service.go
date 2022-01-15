package user

import (
	"context"

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
		router.POST("/user/login", handler.Login)
		return nil
	}})
}

func NewUserService(
	router gin.IRoutes,
	dbConnection db.DatabaseInterface,
	loggerObj logger.LoggerInterface,
	validatorObj validator.ValidatorInterface,
	jwtObj auth.JWTInterface,
) UserServiceInterface {
	fx.New(
		fx.Provide(func() gin.IRoutes { return router }),
		fx.Provide(func() db.DatabaseInterface { return dbConnection }),
		fx.Provide(func() logger.LoggerInterface { return loggerObj }),
		fx.Provide(func() validator.ValidatorInterface { return validatorObj }),
		fx.Provide(func() auth.JWTInterface { return jwtObj }),
		fx.Provide(store.NewStore),
		fx.Provide(usecase.NewUserUseCase),
		fx.Provide(handler.NewHandler),
		fx.Invoke(runAddRoutes),
	).Start(context.TODO())

	return &UserService{
		router: router,
	}
}
