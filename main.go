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
	userGroup := router.Group("/user")
	flashcardGroup := router.Group("/flashcards").Use(middleware.Auth(config.GetSecretKey()))

	auth.NewAuthService(userGroup, dbConnection, loggerObj, config.GetSecretKey())
	flashcard.NewFlashCardService(flashcardGroup, dbConnection, loggerObj)

	router.Run(":" + config.GetPort())

}
