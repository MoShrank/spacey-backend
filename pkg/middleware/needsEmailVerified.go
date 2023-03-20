package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/services/api/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBInterface interface {
	QueryDocument(string, interface{}) *mongo.SingleResult
}

func NeedsEmailVerified(
	store store.StoreInterface,
	jwt auth.JWTInterface,
	cfg config.ConfigInterface,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		verified := c.Query("emailValidated")
		if verified == "true" {
			c.Next()
		} else {
			userID := c.Query("userID")

			user, err := store.GetUserByID(userID)
			if err != nil {
				httpconst.WriteUnauthorized(c, "Invalid user.")
				c.Abort()
			}

			if user.EmailValidated {
				token, err := jwt.CreateJWTWithClaims(userID, map[string]interface{}{
					"IsBeta":         user.BetaUser,
					"EmailValidated": true,
				})
				if err != nil {
					httpconst.WriteUnauthorized(c, "Invalid user.")
					c.Abort()
				}

				c.SetCookie(
					"Authorization",
					token,
					cfg.GetMaxAgeAuth(),
					"/",
					cfg.GetDomain(),
					false,
					true,
				)

				c.Next()
				return
			}

			httpconst.WriteUnauthorized(c, "you need to verify your email before you use spacey.")
			c.Abort()
		}
	}
}
