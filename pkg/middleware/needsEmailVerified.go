package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func NeedsEmailVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		verified := c.Query("emailValidated")
		if verified == "true" {
			c.Next()
		} else {
			httpconst.WriteUnauthorized(c, "you need to verify your email before you use spacey.")
			c.Abort()
		}
	}
}
