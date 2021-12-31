package handler

import (
	"net/http"

	"github.com/moshrank/spacey-backend/pkg/httperror"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"

	"github.com/moshrank/spacey-backend/services/user-service/models"
	"github.com/moshrank/spacey-backend/services/user-service/store"
	"github.com/moshrank/spacey-backend/services/user-service/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	database    store.StoreInterface
	logger      logger.LoggerInterface
	userUsecase usecase.UserUsecaseInterface
	validator   validator.ValidatorInterface
}

type HandlerInterface interface {
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
}

func NewHandler(
	database store.StoreInterface,
	logger logger.LoggerInterface,
	usecase usecase.UserUsecaseInterface,
	validatorObj validator.ValidatorInterface,
) HandlerInterface {
	return &Handler{
		database:    database,
		logger:      logger,
		userUsecase: usecase,
		validator:   validatorObj,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User

	if err := h.validator.ValidateJSON(c, &user); err != nil {
		return
	}

	if user.Name == "" {
		httperror.BadRequest(c)
		return
	}

	user.Password, _ = h.userUsecase.HashPassword(user.Password)

	if err := h.database.SaveUser(&user); err != nil {
		httperror.DatabaseError(c)
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

	if err := h.validator.ValidateJSON(c, &user); err != nil {
		return
	}

	userFromDB, err := h.database.GetUserByEmail(user.Email)

	if err != nil {
		httperror.DatabaseError(c)
		return
	}

	if !h.userUsecase.CheckPasswordHash(user.Password, userFromDB.Password) {
		httperror.Unauthorized(c)
		return
	}

	tokenString, _ := h.userUsecase.CreateJWTWithClaims(userFromDB.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})

}
