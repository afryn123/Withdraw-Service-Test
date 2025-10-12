package middlewares

import (
	"afryn123/withdraw-service/src/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthProtected() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header missing", "No token provided")
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header format", "Invalid token format")
            c.Abort()
            return
        }

        userID, err := utils.ParseJWT(parts[1])
        if err != nil {
            utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
            c.Abort()
            return
        }

        c.Set("userID", userID)
        c.Next()
    }
}

