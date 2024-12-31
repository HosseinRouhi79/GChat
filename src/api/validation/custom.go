package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
	Message  string `json:"message"`
}

func GetValidationError(err error) *[]ValidationError {
	var vErrors []ValidationError
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			var customError ValidationError
			customError.Tag = err.Tag()
			customError.Property = err.Field()
			customError.Value = err.Value().(string)
			vErrors = append(vErrors, customError)
		}
	}
	return &vErrors
}
