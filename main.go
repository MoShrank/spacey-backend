package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service"
	userService "github.com/moshrank/spacey-backend/services/user-service"
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
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.JSONMiddleware())

	router.GET("/ping", ping)

	userGroup := router.Group("/user")
	flashcardGroup := router.Group("/flashcards").Use(middleware.Auth(config.GetSecretKey()))

	userService.NewUserService(
		userGroup,
		dbConnection.GetDB(config.GetUserSeviceDBNAME()),
		loggerObj,
		config.GetSecretKey(),
	)
	flashcard.NewFlashCardService(
		flashcardGroup,
		dbConnection.GetDB(config.GetFlashcardServiceDBName()),
		loggerObj,
	)

	router.Run(":" + config.GetPort())

}
