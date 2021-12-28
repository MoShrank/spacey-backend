package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func validateJWT(tokenString string, secretKey []byte) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	} else {
		return false, err
	}

}

func extractClaims(tokenStr string, secretKey []byte) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing token",
			})
			c.Abort()
			return
		}
		tokenString = tokenString[7:]

		if ok, _ := validateJWT(tokenString, []byte(secretKey)); ok {

			if claims, ok := extractClaims(tokenString, []byte(secretKey)); ok {
				userID := claims["sub"].(string)

				c.Set("userID", userID)

				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid credentials",
				})
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
