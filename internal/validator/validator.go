package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps the go-playground validator
type CustomValidator struct {
    validator *validator.Validate
}

// New returns an echo.Validator
func New() echo.Validator {
    return &CustomValidator{validator: validator.New()}
}

type Normalizable interface {
  Normalize()
}

// Validate satisfies echo.Validator
func (cv *CustomValidator) Validate(i interface{}) error {
    if norm, ok := i.(Normalizable); ok {
    norm.Normalize()
    }

    return cv.validator.Struct(i)
}