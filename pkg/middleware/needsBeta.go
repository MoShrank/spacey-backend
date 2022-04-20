package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func NeedsBeta() gin.HandlerFunc {
	return func(c *gin.Context) {
		isBeta := c.Query("isBeta")
		if isBeta == "true" {
			c.Next()
		} else {
			httpconst.WriteUnauthorized(c, "you need to be part of the closed beta to use this feature!")
			c.Abort()
		}
	}
}
