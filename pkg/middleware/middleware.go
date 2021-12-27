package middleware

import "github.com/gin-gonic/gin"

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.Request.Header.Get("Content-Type")

		if contentType != "application/json" {
			c.JSON(406, gin.H{
				"message": "Content-Type must be application/json",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
