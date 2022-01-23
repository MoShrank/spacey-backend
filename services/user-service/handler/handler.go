package handler

import (
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
}

type HandlerInterface interface {
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

func NewHandler(
	logger logger.LoggerInterface,
	usecase entity.UserUsecaseInterface,
	validatorObj validator.ValidatorInterface,
) HandlerInterface {
	return &Handler{
		logger:      logger,
		userUsecase: usecase,
		validator:   validatorObj,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user entity.UserReq

	if err := h.validator.ValidateJSON(c, &user); err != nil {
		return
	}

	if user.Name == "" {
		httpconst.WriteBadRequest(c)
		return
	}

	_, err := h.userUsecase.CreateUser(&user)

	// TODO create custom errors
	if err != nil {
		httpconst.WriteInternalServerError(c)
		return
	}

	userRes, err := h.userUsecase.Login(user.Email, user.Password)

	if err != nil {
		httpconst.WriteInternalServerError(c)
		return
	}

	// TODO should be set to a secure cookie + expire time should be equal to jwt token expire time
	c.SetCookie("Authorization", userRes.Token, 604800, "/", "", false, true)

	httpconst.WriteCreated(c, userRes)
}

func (h *Handler) Login(c *gin.Context) {
	var user entity.UserReq

	if err := h.validator.ValidateJSON(c, &user); err != nil {
		return
	}

	userRes, err := h.userUsecase.Login(user.Email, user.Password)
	if err != nil {
		httpconst.WriteInternalServerError(c)
		return
	}

	// TODO should be set to a secure cookie + expire time should be equal to jwt token expire time
	c.SetCookie("Authorization", userRes.Token, 604800, "/", "", false, true)

	httpconst.WriteSuccess(c, userRes)

}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	httpconst.WriteSuccess(c, nil)
}
