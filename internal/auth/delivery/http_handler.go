package delivery

import (
	"context"
	"net/http"

	"github.com/afryn123/withdraw-service-test/internal/shared/httpresponse"
	"github.com/afryn123/withdraw-service-test/internal/shared/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service, validate: validator.New()}
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *Handler) Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, "Invalid request payload", "Invalid request payload")
		return
	}
	if err := h.validate.Struct(request); err != nil {
		message := validation.Message(err)
		httpresponse.Error(c, http.StatusBadRequest, message, message)
		return
	}
	token, err := h.service.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		_ = c.Error(err)
		httpresponse.Error(c, http.StatusUnauthorized, "Authentication failed", err.Error())
		return
	}
	httpresponse.Success(c, http.StatusOK, "Login successful", gin.H{"token": token})
}
