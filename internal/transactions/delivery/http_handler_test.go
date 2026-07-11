package delivery_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/mocks"
	"github.com/afryn123/withdraw-service-test/internal/transactions/application"
	"github.com/afryn123/withdraw-service-test/internal/transactions/delivery"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWithdrawHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("returns withdraw result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockTransactionDeliveryService(ctrl)
		userID := uuid.New()
		service.EXPECT().Withdraw(gomock.Any(), userID, int64(40), nil).Return(application.WithdrawResult{
			TransactionID: "transaction-id", Transaction: application.Transaction{BalanceNow: 60},
		}, nil)
		router := transactionRouter(service, userID)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/transaction/withdraw", strings.NewReader(`{"amount":40}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Contains(t, recorder.Body.String(), `"balance_now":60`)
	})

	t.Run("rejects non-positive amount", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		router := transactionRouter(mocks.NewMockTransactionDeliveryService(ctrl), uuid.New())
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/transaction/withdraw", strings.NewReader(`{"amount":0}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

func transactionRouter(service *mocks.MockTransactionDeliveryService, userID uuid.UUID) *gin.Engine {
	router := gin.New()
	group := router.Group("/api/transaction", func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	delivery.RegisterRoutes(group, delivery.NewHandler(service))
	return router
}
