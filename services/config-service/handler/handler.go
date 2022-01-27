package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/config-service/store"
)

type ConfigHandlerInterface interface {
	GetConfig(c *gin.Context)
}

type ConfigHandler struct {
	configStore store.ConfigStoreInterface
	logger      logger.LoggerInterface
}

func NewConfigHandler(
	configStore store.ConfigStoreInterface,
	logger logger.LoggerInterface,
) ConfigHandlerInterface {
	return &ConfigHandler{
		configStore: configStore,
		logger:      logger,
	}
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	configName := c.Param("configName")
	if configName == "" {
		httpconst.WriteBadRequest(c, "configName is required")
		return
	}

	config, err := h.configStore.GetConfig(configName)
	if err != nil {
		httpconst.WriteNotFound(c, err.Error())
		return
	}

	delete(config, "_id")

	httpconst.WriteSuccess(c, config)
}
