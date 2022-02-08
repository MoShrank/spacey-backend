package middleware

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		validHost := false
		var origin *url.URL

		validCORS := []struct {
			Host     string
			Protocol string
		}{
			{
				"localhost:3000",
				"http",
			},
			{
				"spacey.moritz.dev",
				"https",
			},
			{
				"www.spacey.moritz.dev",
				"https",
			},
			{
				"www.spacey-learn.com",
				"https",
			},
			{
				"spacey-learn.com",
				"https",
			},
		}

		referer := c.Request.Header.Get("Referer")
		remote, _ := url.Parse(referer)

		for _, setting := range validCORS {
			if remote.Host == setting.Host && remote.Scheme == setting.Protocol {
				validHost = true
				origin = remote
			}
		}

		if validHost {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin.Scheme+"://"+origin.Host)
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
			c.AbortWithStatus(403)
			return
		}
	}
}
