package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func Auth(authObj auth.JWTInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Request.Cookie("Authorization")

		if err != nil {
			httpconst.WriteUnauthorized(c)
			c.Abort()
			return
		}

		tokenString := authCookie.Value

		if tokenString == "" {
			httpconst.WriteBadRequest(c)
			c.Abort()
			return
		}

		if claims, err := authObj.ValidateJWT(tokenString); err == nil {

			userID := claims.Id
			c.Request.Header.Add("userID", userID)

			c.Next()

		} else {
			httpconst.WriteUnauthorized(c)
			c.Abort()
		}

	}
}
