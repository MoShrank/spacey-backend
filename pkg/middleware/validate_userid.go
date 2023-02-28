package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func ValidateUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Request.URL.Query().Get("userID")

		if userID == "" {
			httpconst.WriteUnauthorized(c, "missing user id")
			return
		}
	}
}
