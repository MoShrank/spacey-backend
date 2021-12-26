package main

import (
	"context"
	"log"
	"os"
	"spacey/auth-service/db"
	"spacey/auth-service/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func runHttpServer(lifecycle fx.Lifecycle, handler handler.HandlerInterface) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		router := gin.Default()

		return router.Run(":8080")
	}})
}

func GetConnectionString() string {
	connectionString := os.Getenv("MONGODB_CONNECTION")
	if connectionString == "" {
		connectionString = "mongodb://192.168.0.204:27017"
	}

	return connectionString
}

func main() {

	ctx := context.TODO()

	app := fx.New(
		fx.Provide(GetConnectionString),
		fx.Provide(db.NewDB),
		fx.Provide(handler.NewHandler),
		fx.Invoke(runHttpServer),
	)
	if err := app.Start(ctx); err != nil {
		log.Println(err)
	}

}
