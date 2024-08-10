package middleware

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func CORSMiddleware(rootDomain string) gin.HandlerFunc {
	return func(c *gin.Context) {

		validHost := false

		origin := c.Request.Header.Get("Origin")

		originParsed, _ := url.Parse(origin)
		originHost := originParsed.Host
		originHostSplit := strings.Split(originHost, ":")
		originHostCleaned := originHostSplit[0]

		if origin == "" {
			c.Next()
			return
		}

		if originHostCleaned == rootDomain {
			validHost = true
		} else if originHostCleaned == "localhost" {
			validHost = true
		} else {
			validHost = false
		}

		if validHost {
			c.Writer.Header().
				Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().
				Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		} else {
			httpconst.WriteUnauthorized(c, "cors error. invalid origin: "+origin)
			return
		}
	}
}
