package delivery_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/auth/delivery"
	"github.com/afryn123/withdraw-service-test/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("rejects missing header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		router := middlewareRouter(mocks.NewMockTokenManager(ctrl))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/protected", nil))
		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tokens := mocks.NewMockTokenManager(ctrl)
		tokens.EXPECT().Parse("bad").Return(uuid.Nil, errors.New("invalid"))
		router := middlewareRouter(tokens)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/protected", nil)
		request.Header.Set("Authorization", "Bearer bad")
		router.ServeHTTP(recorder, request)
		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("passes authenticated user id to next handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tokens := mocks.NewMockTokenManager(ctrl)
		userID := uuid.New()
		tokens.EXPECT().Parse("valid").Return(userID, nil)
		router := middlewareRouter(tokens)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/protected", nil)
		request.Header.Set("Authorization", "Bearer valid")
		router.ServeHTTP(recorder, request)
		require.Equal(t, http.StatusOK, recorder.Code)
		require.Equal(t, userID.String(), recorder.Body.String())
	})
}

func middlewareRouter(tokens *mocks.MockTokenManager) *gin.Engine {
	router := gin.New()
	router.GET("/protected", delivery.Middleware(tokens), func(c *gin.Context) {
		c.String(http.StatusOK, c.MustGet("userID").(uuid.UUID).String())
	})
	return router
}
