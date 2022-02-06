package validator

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

type Validator struct {
	validator *validator.Validate
}

type ValidatorInterface interface {
	Validate(v interface{}) error
	ValidateJSON(c *gin.Context, obj interface{}) error
}

type ValidationError struct {
	Field string `json:"field"`
	Msg   string `json:"message"`
}

func NewValidator() ValidatorInterface {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(obj interface{}) error {
	return v.validator.Struct(obj)
}

func (v *Validator) msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}

func (v *Validator) ValidateJSON(c *gin.Context, obj interface{}) error {

	err := c.ShouldBind(&obj)

	fmt.Println(obj)
	fmt.Println("err", err)
	if err != nil {

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ValidationError, len(ve))
			for i, fe := range ve {
				out[i] = ValidationError{fe.Field(), v.msgForTag(fe.Tag())}
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   httpconst.ErrorMapping[http.StatusBadRequest],
				"message": out,
			})

		} else {
			httpconst.WriteValidationError(c, err.Error())
		}
		return err
	}

	return nil
}
