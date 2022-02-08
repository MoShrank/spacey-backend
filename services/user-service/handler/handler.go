package handler

import (
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"

	"github.com/moshrank/spacey-backend/services/user-service/entity"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger      logger.LoggerInterface
	userUsecase entity.UserUsecaseInterface
	validator   validator.ValidatorInterface
	config      config.ConfigInterface
}

type HandlerInterface interface {
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetUser(c *gin.Context)
}

func NewHandler(
	logger logger.LoggerInterface,
	usecase entity.UserUsecaseInterface,
	validatorObj validator.ValidatorInterface,
	configObj config.ConfigInterface,
) HandlerInterface {
	return &Handler{
		logger:      logger,
		userUsecase: usecase,
		validator:   validatorObj,
		config:      configObj,
	}
}

func (h *Handler) setAuthCookie(c *gin.Context, token string) {
	c.SetCookie(
		"Authorization",
		token,
		h.config.GetMaxAgeAuth(),
		"/",
		h.config.GetDomain(),
		false,
		true,
	)
	c.SetCookie(
		"LoggedIn",
		"true",
		h.config.GetMaxAgeAuth(),
		"/",
		h.config.GetDomain(),
		false,
		false,
	)

}

func (h *Handler) CreateUser(c *gin.Context) {
	var user entity.UserReq

	if err := h.validator.ValidateJSON(c, &user); err != nil {
		return
	}

	if user.Name == "" {
		httpconst.WriteBadRequest(c, "name is required.")
		return
	}

	_, err := h.userUsecase.CreateUser(&user)

	if err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	userRes, err := h.userUsecase.Login(user.Email, user.Password)

	if err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	h.setAuthCookie(c, userRes.Token)

	httpconst.WriteCreated(c, userRes)
}

func (h *Handler) Login(c *gin.Context) {
	var user entity.UserReq

	if err := h.validator.ValidateJSON(c, &user); err != nil {
		return
	}

	userRes, err := h.userUsecase.Login(user.Email, user.Password)
	if err != nil {
		httpconst.WriteBadRequest(c, "Invalid email or password.")
		return
	}

	h.setAuthCookie(c, userRes.Token)

	httpconst.WriteSuccess(c, userRes)

}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", h.config.GetDomain(), false, true)
	c.SetCookie("LoggedIn", "false", -1, "/", h.config.GetDomain(), false, false)
	httpconst.WriteSuccess(c, nil)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteBadRequest(c, "userID is required.")
		return
	}

	userRes, err := h.userUsecase.GetUserByID(userID)
	if err != nil {
		httpconst.WriteNotFound(c, "Could not find userID in database.")
		return
	}

	httpconst.WriteSuccess(c, userRes)
}
