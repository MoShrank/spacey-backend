package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		validHost := false

		validDomains := []string{
			"http://localhost:3000",
			"https://www.spacey.moritz.dev",
			"https://spacey.moritz.dev",
		}

		referer := c.Request.Header.Get("Referer")

		for _, domain := range validDomains {
			if referer == domain {
				validHost = true
			}
		}

		if validHost {
			c.Writer.Header().Set("Access-Control-Allow-Origin", referer)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().
				Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "invalid host",
			})
			return
		}
	}
}
