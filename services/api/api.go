package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/api/routes"
)

func main() {
	config, err := config.NewConfig()
	log := logger.NewLogger(config)

	router := gin.New()
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery())

	routes.CreateRoutes(router, config)

	if err != nil {
		panic(err)
	}
	log.Info("Starting server on port: " + config.GetPort())
	router.Run(":" + config.GetPort())
}
