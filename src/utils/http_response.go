package utils

import (
	"afryn123/withdraw-service/src/dtos"

	"github.com/gin-gonic/gin"
)


func SuccessResponse(c *gin.Context, status int, message string, data any) {
	c.JSON(status, dtos.APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
		Error:   nil,	
	})
}

func ErrorResponse(c *gin.Context, status int, message string, error interface{}) {
	c.JSON(status, dtos.APIResponse{
		Status:  false,
		Message: message,
		Error:   error,
	})
}

func ValidationErrorResponse(err error) string {
	message := FormatValidationError(err)
	return message
}
