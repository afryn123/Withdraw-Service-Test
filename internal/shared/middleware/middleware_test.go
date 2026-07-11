package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sharedmiddleware "github.com/afryn123/withdraw-service-test/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	router := gin.New()
	router.Use(sharedmiddleware.RequestLogger(logger), sharedmiddleware.Recovery(logger))
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("X-Request-ID", "request-123")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, "request-123", recorder.Header().Get("X-Request-ID"))
	logLine := output.String()
	require.Contains(t, logLine, `"msg":"http request completed"`)
	require.Contains(t, logLine, `"request_id":"request-123"`)
	require.Contains(t, logLine, `"method":"GET"`)
	require.Contains(t, logLine, `"path":"/health"`)
	require.Contains(t, logLine, `"status":204`)
}

func TestRecoveryUsesStructuredLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	router := gin.New()
	router.Use(sharedmiddleware.RequestLogger(logger), sharedmiddleware.Recovery(logger))
	router.GET("/panic", func(*gin.Context) { panic("unexpected") })
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/panic", nil))

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	logOutput := output.String()
	require.Contains(t, logOutput, `"msg":"panic recovered"`)
	require.Contains(t, logOutput, `"panic":"unexpected"`)
	require.True(t, strings.Count(logOutput, "\n") >= 2, "panic and request completion must both be logged")
}
