package httperror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
}

func DatabaseError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Database error",
	})
}

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Bad request",
	})
}

func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Validation error",
		"message": err,
	})
}
