package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func validateJWT(tokenString string, secretKey string) (bool, error) {
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

func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authentication")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing token",
			})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]

		if ok, _ := validateJWT(tokenString, secretKey); ok {
			c.JSON(http.StatusOK, gin.H{
				"message": "Authentication successful",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
		}
	}
}
