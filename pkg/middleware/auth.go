package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func Auth(authObj auth.JWTInterface, config config.ConfigInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Request.Cookie("Authorization")
		if err != nil {
			httpconst.WriteUnauthorized(c, "missing authorization cookie")
			c.Abort()
			return
		}

		tokenString := authCookie.Value

		if tokenString == "" {
			httpconst.WriteBadRequest(c, "could not find authorization token.")
			c.Abort()
			return
		}

		if claims, err := authObj.ValidateJWT(tokenString, []string{"IsBeta", "EmailValidated"}); err == nil {
			userID := claims["Id"].(string)
			isBeta := claims["IsBeta"].(bool)
			emailValidated := claims["EmailValidated"].(bool)

			q := c.Request.URL.Query()
			q.Del("userID")
			q.Add("userID", userID)
			q.Del("isBeta")
			q.Add("isBeta", fmt.Sprintf("%t", isBeta))
			q.Del("emailValidated")
			q.Add("emailValidated", fmt.Sprintf("%t", emailValidated))
			c.Request.URL.RawQuery = q.Encode()

			c.Next()

		} else {
			c.SetCookie("Authorization", "", -1, "/", config.GetDomain(), false, true)
			c.SetCookie("LoggedIn", "false", -1, "/", config.GetDomain(), false, false)
			httpconst.WriteUnauthorized(c, "invalid authorization token")
			c.Abort()
		}

	}
}
