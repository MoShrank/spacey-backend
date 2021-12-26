package auth
//hello
import (
	"context"
	"spacey/auth-service/handler"
	"spacey/auth-service/store"

	"spacey/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type AuthService struct {
	router  *gin.RouterGroup
	db      *db.Database
	handler *handler.Handler
}

type AuthServiceInterface interface {
}

func runAddRoutes(
	lifecycle fx.Lifecycle,
	handler handler.HandlerInterface,
	router *gin.RouterGroup,
	db *db.Database,
) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router.GET("/", handler.Ping)
		router.POST("/user", handler.CreateUser)
		router.GET("/user/login", handler.Login)
		router.GET("/user/auth", handler.ValidateJWT)
		return nil
	}})
}
func NewAuthService(router *gin.RouterGroup, db *db.Database) AuthServiceInterface {
	fx.New(
		fx.Provide(store.NewStore),
		fx.Provide(handler.NewHandler),
		fx.Invoke(runAddRoutes),
	).Run()

	return &AuthService{
		router: router,
	}
}
