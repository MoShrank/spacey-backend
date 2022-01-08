package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service"
	"github.com/moshrank/spacey-backend/services/user-service"
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

	loggerObj := logger.NewLogger(config.GetLogLevel(), config.GetGrayLogConnection(), config)
	dbConnection := db.NewDB(config.GetMongoDBConnection(), loggerObj)
	validator := validator.NewValidator()
	authObj := auth.NewJWT(config.GetJWTSecret())

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.JSONMiddleware())

	router.GET("/ping", ping)

	flashcardGroup := router.Group("/flashcards").Use(middleware.Auth(authObj))

	user.NewUserService(
		router,
		dbConnection.GetDB(config.GetUserSeviceDBNAME()),
		loggerObj,
		validator,
		authObj,
	)
	flashcard.NewFlashCardService(
		flashcardGroup,
		dbConnection.GetDB(config.GetFlashcardServiceDBName()),
		loggerObj,
		validator,
	)

	router.Run(":" + config.GetPort())

}
