package middlewares

import (
	"afryn123/withdraw-service/src/utils"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func CustomRecoverPanic() gin.HandlerFunc {
	logger := utils.GetLogger("system")
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := debug.Stack()
				logger.Printf("[ERROR][PANIC] %v\n%s\n", r, stackTrace)

				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", "Server encountered an unexpected condition")
				c.Abort()
			}
		}()
		c.Next()
	}
}
