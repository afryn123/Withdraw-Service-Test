package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/afryn123/withdraw-service-test/internal/shared/httpresponse"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "requestID"

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE", "GET", "OPTIONS"},
		AllowHeaders: []string{"*"},
	})
}

func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()

		attributes := []any{
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(startedAt).Milliseconds(),
			"client_ip", c.ClientIP(),
		}
		if len(c.Errors) > 0 {
			attributes = append(attributes, "error", c.Errors.String())
		}
		switch status := c.Writer.Status(); {
		case status >= http.StatusInternalServerError:
			logger.ErrorContext(c.Request.Context(), "http request completed", attributes...)
		case status >= http.StatusBadRequest:
			logger.WarnContext(c.Request.Context(), "http request completed", attributes...)
		default:
			logger.InfoContext(c.Request.Context(), "http request completed", attributes...)
		}
	}
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.ErrorContext(c.Request.Context(), "panic recovered",
					"request_id", c.GetString(RequestIDKey),
					"panic", recovered,
					"stack", string(debug.Stack()),
				)
				httpresponse.Error(c, http.StatusInternalServerError, "Internal Server Error", "Server encountered an unexpected condition")
				c.Abort()
			}
		}()
		c.Next()
	}
}
