package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/auth"
	"github.com/moshrank/spacey-backend/services/flashcard-management-service"
)

func GetConnectionString() string {
	connectionString := os.Getenv("MONGODB_CONNECTION")
	if connectionString == "" {
		connectionString = "mongodb://localhost:27017"
	}

	return connectionString
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return port
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	PORT := GetPort()

	connectionString := GetConnectionString()
	loggerObj := logger.NewLogger()
	dbConnection := db.NewDB(connectionString, loggerObj)

	router := gin.Default()
	router.Use(middleware.JSONMiddleware())

	router.GET("/ping", ping)
	authGroup := router.Group("/auth")
	flashcardGroup := router.Group("/flashcards")

	auth.NewAuthService(authGroup, dbConnection, loggerObj)
	flashcard.NewFlashCardService(flashcardGroup, dbConnection, loggerObj)

	router.Run(":" + PORT)

}
