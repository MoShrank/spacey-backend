package httpconst

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Created",
		"data":    data,
	})
}

func WriteSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    data,
	})
}
