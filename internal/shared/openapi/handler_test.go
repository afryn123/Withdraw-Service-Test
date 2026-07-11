package openapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/shared/openapi"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestOpenAPIRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	openapi.RegisterRoutes(router, openapi.NewHandler())

	t.Run("serves a valid embedded specification", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil))

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Contains(t, recorder.Header().Get("Content-Type"), "application/yaml")
		var specification map[string]any
		require.NoError(t, yaml.Unmarshal(recorder.Body.Bytes(), &specification))
		require.Equal(t, "3.0.3", specification["openapi"])
		require.Contains(t, specification, "paths")
	})

	t.Run("redirects to the trailing-slash documentation URL", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/docs", nil))

		require.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
		require.Equal(t, "/docs/", recorder.Header().Get("Location"))
	})

	t.Run("uses a relative specification URL", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/docs/", nil))

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Contains(t, recorder.Body.String(), `url: "../openapi.yaml"`)
	})
}
