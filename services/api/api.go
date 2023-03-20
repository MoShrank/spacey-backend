package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/api/routes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {

	config, err := config.NewConfig()
	log := logger.NewLogger(config)
	db := db.NewDB(config, log)

	router := gin.New()
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery())
	router.Use(middleware.PrometheusMiddleware())

	router.GET("/prometheus", prometheusHandler())

	routes.CreateRoutes(router, config, db)

	if err != nil {
		panic(err)
	}
	log.Info("Starting server on port: " + config.GetPort())
	router.Run(":" + config.GetPort())
}
