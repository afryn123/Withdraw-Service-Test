package delivery

import (
	"context"
	"net/http"

	"github.com/afryn123/withdraw-service-test/internal/shared/httpresponse"
	"github.com/afryn123/withdraw-service-test/internal/wallet/application"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	FindBalanceByUserID(ctx context.Context, userID uuid.UUID) (application.Balance, error)
}

type Handler struct{ service Service }

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) GetBalance(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, "Invalid user ID format", "Invalid user ID format")
		return
	}
	balance, err := h.service.FindBalanceByUserID(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		httpresponse.Error(c, http.StatusInternalServerError, "Failed to retrieve balance", "Failed to retrieve balance")
		return
	}
	httpresponse.Success(c, http.StatusOK, "Balance retrieved successfully", balance)
}
