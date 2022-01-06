package middleware

import (
	"net/http"

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
			httpconst.WriteUnauthorized(c)
			c.Abort()
			return
		}
		tokenString = tokenString[7:]

		if ok, _ := authObj.ValidateJWT(tokenString); ok {

			if claims, ok := authObj.ExtractClaims(tokenString); ok {
				userID := claims["sub"].(string)

				c.Set("userID", userID)

				c.Next()
			} else {
				httpconst.WriteUnauthorized(c)
				c.Abort()
				return
			}

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token!",
			})
			c.Abort()
		}

	}
}
