package delivery_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/mocks"
	"github.com/afryn123/withdraw-service-test/internal/wallet/application"
	"github.com/afryn123/withdraw-service-test/internal/wallet/delivery"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetBalanceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("returns balance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockWalletDeliveryService(ctrl)
		userID, walletID := uuid.New(), uuid.New()
		service.EXPECT().FindBalanceByUserID(gomock.Any(), userID).Return(application.Balance{WalletID: walletID, Balance: 100}, nil)
		router := gin.New()
		delivery.RegisterRoutes(router.Group("/api/wallet"), delivery.NewHandler(service))
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/wallet/"+userID.String()+"/balance", nil))

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Contains(t, recorder.Body.String(), `"balance":100`)
	})

	t.Run("rejects malformed user id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		router := gin.New()
		delivery.RegisterRoutes(router.Group("/api/wallet"), delivery.NewHandler(mocks.NewMockWalletDeliveryService(ctrl)))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/wallet/bad/balance", nil))
		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
