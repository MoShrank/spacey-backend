package validator

import (
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

func NewValidator() ValidatorInterface {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(obj interface{}) error {
	return v.validator.Struct(obj)
}

func (v *Validator) ValidateJSON(c *gin.Context, obj interface{}) error {
	err := c.BindJSON(obj)
	if err != nil {
		httpconst.WriteValidationError(c, err.Error())
		return err
	}

	if err := v.Validate(obj); err != nil {
		httpconst.WriteValidationError(c, err.Error())
		return err
	}

	return nil
}
