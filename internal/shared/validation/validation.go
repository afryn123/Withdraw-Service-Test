package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func Message(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return "invalid request"
	}
	messages := make([]string, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		field := fieldError.Field()
		switch fieldError.Tag() {
		case "required":
			messages = append(messages, field+" is required")
		case "email":
			messages = append(messages, field+" must be a valid email")
		case "min":
			messages = append(messages, field+" must be at least "+fieldError.Param()+" characters")
		case "gt":
			messages = append(messages, field+" must be greater than "+fieldError.Param())
		default:
			messages = append(messages, field+" is invalid")
		}
	}
	return strings.Join(messages, ", ")
}
