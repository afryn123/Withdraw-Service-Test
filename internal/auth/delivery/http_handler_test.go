package delivery_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/auth/delivery"
	authdomain "github.com/afryn123/withdraw-service-test/internal/auth/domain"
	"github.com/afryn123/withdraw-service-test/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("returns token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockAuthDeliveryService(ctrl)
		service.EXPECT().Login(gomock.Any(), "a@example.com", "secret12").Return("jwt", nil)
		router := gin.New()
		delivery.RegisterRoutes(router.Group("/api/auth"), delivery.NewHandler(service))
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"a@example.com","password":"secret12"}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Contains(t, recorder.Body.String(), `"token":"jwt"`)
	})

	t.Run("maps invalid credentials to unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockAuthDeliveryService(ctrl)
		service.EXPECT().Login(gomock.Any(), "a@example.com", "secret12").Return("", authdomain.ErrInvalidCredentials)
		router := gin.New()
		delivery.RegisterRoutes(router.Group("/api/auth"), delivery.NewHandler(service))
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"a@example.com","password":"secret12"}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("rejects invalid email before calling service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		router := gin.New()
		delivery.RegisterRoutes(router.Group("/api/auth"), delivery.NewHandler(mocks.NewMockAuthDeliveryService(ctrl)))
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"bad","password":"secret12"}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
