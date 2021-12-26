package main

import (
	"spacey/auth"
	"os"

	"github.com/gin-gonic/gin"
)

func GetConnectionString() string {
	connectionString := os.Getenv("MONGODB_CONNECTION")
	if connectionString == "" {
		connectionString = "mongodb://192.168.0.204:27017"
	}

	return connectionString
}

func main() {
	router := gin.Default()

	authGroup := router.Group("/auth")
	auth.NewAuthService(authGroup)

	router.Run(":8080")
}
