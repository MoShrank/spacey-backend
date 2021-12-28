package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/moshrank/spacey-backend/pkg/logger"

	"github.com/moshrank/spacey-backend/services/auth/models"
	"github.com/moshrank/spacey-backend/services/auth/store"
	"github.com/moshrank/spacey-backend/services/auth/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	database    store.StoreInterface
	logger      logger.LoggerInterface
	userUsecase usecase.UserUsecaseInterface
	validator   *validator.Validate
}

type HandlerInterface interface {
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
	Authenticate(c *gin.Context)
}

func NewHandler(
	database store.StoreInterface,
	logger logger.LoggerInterface,
	usecase usecase.UserUsecaseInterface,
) HandlerInterface {
	return &Handler{
		database:    database,
		logger:      logger,
		userUsecase: usecase,
		validator:   validator.New(),
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)

	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.validator.Struct(user)

	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Password, _ = h.userUsecase.HashPassword(user.Password)

	err = h.database.SaveUser(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userFromDB, err := h.database.GetUserByEmail(user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !h.userUsecase.CheckPasswordHash(user.Password, userFromDB.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	tokenString, _ := h.userUsecase.CreateJWTWithClaims(userFromDB.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})

}

func (h *Handler) Authenticate(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authentication")
	tokenString = tokenString[7:]

	if ok, _ := h.userUsecase.ValidateJWT(tokenString); ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication successful",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
	}

}
