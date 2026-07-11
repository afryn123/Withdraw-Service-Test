package delivery

import (
	"context"
	"errors"
	"net/http"

	"github.com/afryn123/withdraw-service-test/internal/shared/httpresponse"
	"github.com/afryn123/withdraw-service-test/internal/shared/validation"
	"github.com/afryn123/withdraw-service-test/internal/user/application"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, command application.CreateUserCommand) error
	Update(ctx context.Context, command application.UpdateUserCommand) error
	Delete(ctx context.Context, command application.DeleteUserCommand) error
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service, validate: validator.New()}
}

type createRequest struct {
	Name      string  `json:"name" validate:"required"`
	Username  string  `json:"username" validate:"required"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	CreatedBy *string `json:"created_by" validate:"required"`
}

type updateRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=1"`
	Username *string `json:"username" validate:"omitempty,min=1"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=6"`
}

func (h *Handler) Create(c *gin.Context) {
	var request createRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, "Invalid request payload", "Invalid request payload")
		return
	}
	if err := h.validate.Struct(request); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, validation.Message(err), nil)
		return
	}
	err := h.service.Create(c.Request.Context(), application.CreateUserCommand{
		Name: request.Name, Username: request.Username, Email: request.Email,
		Password: request.Password, CreatedBy: request.CreatedBy,
	})
	if err != nil {
		_ = c.Error(err)
		httpresponse.Error(c, http.StatusInternalServerError, "Failed to create user", "Failed to create user")
		return
	}
	httpresponse.Success(c, http.StatusCreated, "User created successfully", nil)
}

func (h *Handler) Update(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, "Invalid user ID format", "Invalid user ID format")
		return
	}
	actorID, ok := c.Get("userID")
	if !ok {
		httpresponse.Error(c, http.StatusUnauthorized, "User not authenticated", "User not authenticated")
		return
	}
	parsedActorID, ok := actorID.(uuid.UUID)
	if !ok {
		httpresponse.Error(c, http.StatusInternalServerError, "Invalid authenticated user ID", nil)
		return
	}

	var request updateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, "Invalid request payload", "Invalid request payload")
		return
	}
	if request.Name == nil && request.Username == nil && request.Email == nil && request.Password == nil {
		httpresponse.Error(c, http.StatusBadRequest, "At least one field must be provided", nil)
		return
	}
	if err := h.validate.Struct(request); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, validation.Message(err), nil)
		return
	}

	err = h.service.Update(c.Request.Context(), application.UpdateUserCommand{
		UserID: userID, ActorID: parsedActorID, Name: request.Name,
		Username: request.Username, Email: request.Email, Password: request.Password,
	})
	if err != nil {
		_ = c.Error(err)
	}
	switch {
	case errors.Is(err, userdomain.ErrForbidden):
		httpresponse.Error(c, http.StatusForbidden, "Forbidden", err.Error())
	case errors.Is(err, userdomain.ErrNotFound):
		httpresponse.Error(c, http.StatusNotFound, "User not found", err.Error())
	case err != nil:
		httpresponse.Error(c, http.StatusInternalServerError, "Failed to update user", "Failed to update user")
	default:
		httpresponse.Success(c, http.StatusOK, "User updated successfully", nil)
	}
}

func (h *Handler) Delete(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, "Invalid user ID format", "Invalid user ID format")
		return
	}
	actorID, ok := c.Get("userID")
	if !ok {
		httpresponse.Error(c, http.StatusUnauthorized, "User not authenticated", "User not authenticated")
		return
	}
	parsedActorID, ok := actorID.(uuid.UUID)
	if !ok {
		httpresponse.Error(c, http.StatusInternalServerError, "Invalid authenticated user ID", nil)
		return
	}

	err = h.service.Delete(c.Request.Context(), application.DeleteUserCommand{
		UserID: userID, ActorID: parsedActorID,
	})
	if err != nil {
		_ = c.Error(err)
	}
	switch {
	case errors.Is(err, userdomain.ErrForbidden):
		httpresponse.Error(c, http.StatusForbidden, "Forbidden", err.Error())
	case errors.Is(err, userdomain.ErrNotFound):
		httpresponse.Error(c, http.StatusNotFound, "User not found", err.Error())
	case err != nil:
		httpresponse.Error(c, http.StatusInternalServerError, "Failed to delete user", "Failed to delete user")
	default:
		httpresponse.Success(c, http.StatusOK, "User deleted successfully", nil)
	}
}
