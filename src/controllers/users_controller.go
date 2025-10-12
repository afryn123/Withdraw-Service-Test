package controllers

import (
	"afryn123/withdraw-service/src/dtos"
	"afryn123/withdraw-service/src/services"
	"afryn123/withdraw-service/src/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UsersContoller struct {
	usersService *services.UsersService
	validator    *validator.Validate
	logger       *log.Logger
}

func NewUsersContoller(usersService *services.UsersService) *UsersContoller {
	return &UsersContoller{
		usersService: usersService,
		validator:    validator.New(),
		logger:       utils.GetLogger("system"),
	}
}

// Create User
func (c *UsersContoller) Create(ctx *gin.Context) {
	dto := new(dtos.UserCreateDTO)
	if err := ctx.BindJSON(dto); err != nil {
		c.logger.Println("[ERROR][CREATE USER] Error binding JSON:", err)
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authorization header missing", "No token provided")
		return
	}

	if err := c.validator.Struct(dto); err != nil {
		message := utils.ValidationErrorResponse(err)
		c.logger.Println("[ERROR][CREATE USER] Validation error:", message)
		utils.ErrorResponse(ctx, http.StatusBadRequest, message, nil)
		return
	}

	if err := c.usersService.Create(dto); err != nil {
		c.logger.Println("[ERROR][CREATE USER] Service error:", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create user", "Failed to create user")
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "User created successfully", nil)
}
