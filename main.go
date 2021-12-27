package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/auth"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service"
)

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {

	config, err := config.NewConfig()

	if err != nil {
		panic(err)
	}

	loggerObj := logger.NewLogger(config.GetLogLevel())
	dbConnection := db.NewDB(config.GetMongoDBConnection(), loggerObj)

	router := gin.Default()
	router.Use(middleware.JSONMiddleware())

	router.GET("/ping", ping)
	authGroup := router.Group("/auth")
	flashcardGroup := router.Group("/flashcards")

	auth.NewAuthService(authGroup, dbConnection, loggerObj)
	flashcard.NewFlashCardService(flashcardGroup, dbConnection, loggerObj)

	router.Run(":" + config.GetPort())

}
