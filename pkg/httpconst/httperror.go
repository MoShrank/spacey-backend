package httpconst

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrorMapping = map[int]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusMethodNotAllowed:    "Method Not Allowed",
	http.StatusInternalServerError: "Internal Server Error",
}

func WriteUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": ErrorMapping[http.StatusUnauthorized],
	})
}

func WriteDatabaseError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": ErrorMapping[http.StatusInternalServerError],
	})
}

func WriteBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":    ErrorMapping[http.StatusBadRequest],
		"messsage": message,
	})
}

func WriteValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   ErrorMapping[http.StatusBadRequest],
		"message": message,
	})
}

func WriteInternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   ErrorMapping[http.StatusInternalServerError],
		"message": message,
	})
}

func WriteNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"error":    ErrorMapping[http.StatusNotFound],
		"messsage": message,
	})
}

func WriteLimitReached(c *gin.Context) {
	c.JSON(429, gin.H{
		"error": "Too many requests",
	})
}
