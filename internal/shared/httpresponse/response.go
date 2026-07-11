package httpresponse

import "github.com/gin-gonic/gin"

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error"`
}

func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Response{Status: true, Message: message, Data: data})
}

func Error(c *gin.Context, status int, message string, detail any) {
	c.JSON(status, Response{Status: false, Message: message, Error: detail})
}
