package delivery

import (
	"context"
	"net/http"

	"github.com/afryn123/withdraw-service-test/internal/shared/httpresponse"
	"github.com/afryn123/withdraw-service-test/internal/shared/validation"
	"github.com/afryn123/withdraw-service-test/internal/transactions/application"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service interface {
	Withdraw(ctx context.Context, userID uuid.UUID, amount int64, remark *string) (application.WithdrawResult, error)
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service, validate: validator.New()}
}

type withdrawRequest struct {
	Amount int64   `json:"amount" validate:"required,gt=0"`
	Remark *string `json:"remark"`
}

func (h *Handler) Withdraw(c *gin.Context) {
	var request withdrawRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, "Invalid request payload", "Invalid request payload")
		return
	}
	if err := h.validate.Struct(request); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, validation.Message(err), nil)
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		httpresponse.Error(c, http.StatusUnauthorized, "User not authenticated", "User not authenticated")
		return
	}
	parsedUserID, ok := userID.(uuid.UUID)
	if !ok {
		httpresponse.Error(c, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}
	result, err := h.service.Withdraw(c.Request.Context(), parsedUserID, request.Amount, request.Remark)
	if err != nil {
		_ = c.Error(err)
		httpresponse.Error(c, http.StatusInternalServerError, "Failed to process withdraw", err.Error())
		return
	}
	httpresponse.Success(c, http.StatusOK, "Withdraw processed successfully", result)
}
