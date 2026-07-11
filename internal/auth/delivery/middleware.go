package delivery

import (
	"net/http"
	"strings"

	authdomain "github.com/afryn123/withdraw-service-test/internal/auth/domain"
	"github.com/afryn123/withdraw-service-test/internal/shared/httpresponse"
	"github.com/gin-gonic/gin"
)

func Middleware(tokens authdomain.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			httpresponse.Error(c, http.StatusUnauthorized, "Authorization header missing", "No token provided")
			c.Abort()
			return
		}
		parts := strings.Fields(header)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpresponse.Error(c, http.StatusUnauthorized, "Invalid Authorization header format", "Invalid token format")
			c.Abort()
			return
		}
		userID, err := tokens.Parse(parts[1])
		if err != nil {
			_ = c.Error(err)
			httpresponse.Error(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
