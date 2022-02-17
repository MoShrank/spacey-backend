package httpconst

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteCreated(c *gin.Context, data interface{}) {
	message := gin.H{
		"message": "Created",
	}

	if data != nil {
		message["data"] = data
	}

	c.JSON(http.StatusCreated, message)
}

func WriteSuccess(c *gin.Context, data interface{}) {
	message := gin.H{
		"message": "Success",
	}

	if data != nil {
		message["data"] = data
	}

	c.JSON(http.StatusOK, message)
}
