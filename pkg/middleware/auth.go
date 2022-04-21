package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func Auth(authObj auth.JWTInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Request.Cookie("Authorization")

		if err != nil {
			httpconst.WriteUnauthorized(c, "missing authorization cookie")
			c.Abort()
			return
		}

		tokenString := authCookie.Value

		if tokenString == "" {
			httpconst.WriteBadRequest(c, "Could not find authorization token.")
			c.Abort()
			return
		}

		if claims, err := authObj.ValidateJWT(tokenString); err == nil {

			userID := claims["Id"].(string)
			isBeta := claims["IsBeta"].(bool)
			//construct new url with userID as query parameter and set it in request
			q := c.Request.URL.Query()
			q.Del("userID")
			q.Add("userID", userID)
			q.Del("isBeta")
			q.Add("isBeta", fmt.Sprintf("%t", isBeta))
			c.Request.URL.RawQuery = q.Encode()

			c.Next()

		} else {
			httpconst.WriteUnauthorized(c, "invalid authorization token")
			c.Abort()
		}

	}
}
