package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	log.Print("Ping")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {

	router := gin.Default()
	router.GET("/", Ping)

	router.Run(":8080")

}
