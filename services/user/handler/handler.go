package handler

import (
	"fmt"
	"net/http"
	"spacey/auth-service/entities"
	"spacey/auth-service/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte = []byte("secret")

type Handler struct {
	database store.StoreInterface
}

type HandlerInterface interface {
	Ping(c *gin.Context)
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
	ValidateJWT(c *gin.Context)
}

func NewHandler(database store.StoreInterface) HandlerInterface {
	return &Handler{
		database: database,
	}
}

func (h *Handler) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *Handler) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Print(err)
	return err == nil
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user entities.User
	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Print(err.Error())
		return
	}

	user.CreatedAt = time.Now()
	user.Password, _ = h.hashPassword(user.Password)

	_, err = h.database.SaveUser(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var user entities.User
	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Print(err.Error())
		return
	}

	password, err := h.database.GetPassword(user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !h.checkPasswordHash(user.Password, password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})

}

func (h *Handler) ValidateJWT(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authentication")
	tokenString = tokenString[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Authentication successful",
		})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

}
