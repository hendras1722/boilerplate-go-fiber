package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a single validation error field and its message
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationError converts validator.ValidationErrors into an array of ValidationError structs
func FormatValidationError(err error) []ValidationError {
	var errs []ValidationError
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, vErr := range validationErrs {
			errs = append(errs, ValidationError{
				Field:   vErr.Field(),
				Message: getValidationMessage(vErr),
			})
		}
	} else {
		errs = append(errs, ValidationError{
			Field:   "unknown",
			Message: err.Error(),
		})
	}
	return errs
}

// getValidationMessage returns a custom error message based on the validator tag
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("Field %s is required", err.Field())
	case "email":
		return fmt.Sprintf("Field %s failed validation on the 'email' tag", err.Field())
	case "min":
		if err.Field() == "Password" {
			return fmt.Sprintf("Password minimal harus %s karakter", err.Param())
		}
		return fmt.Sprintf("Field %s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("Field %s must be at most %s characters", err.Field(), err.Param())
	default:
		return fmt.Sprintf("Field %s failed validation on the '%s' tag", err.Field(), err.Tag())
	}
}
