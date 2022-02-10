package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/api/routes"
)

type API struct {
	router *gin.Engine
}

type APIInterface interface {
	Run(port string)
}

func NewAPI(config config.ConfigInterface) APIInterface {
	logger := logger.NewLogger(config)
	fmt.Println("fmt: Starting API...")
	logger.Info("logger: Starting API...")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.Recovery())

	routes.CreateRoutes(router, config)

	return &API{
		router: router,
	}
}

func (api *API) Run(port string) {
	api.router.Run(":" + port)
}

func main() {
	config, err := config.NewConfig()

	if err != nil {
		panic(err)
	}

	api := NewAPI(config)

	api.Run(config.GetPort())
}
