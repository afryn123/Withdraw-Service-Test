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

type AuthContoller struct {
	authService *services.AuthService
	validator   *validator.Validate
	logger      *log.Logger
}

func NewAuthController(authService *services.AuthService) *AuthContoller {
	return &AuthContoller{
		authService: authService,
		validator:   validator.New(),
		logger:      utils.GetLogger("system"),
	}
}

func (c *AuthContoller) Login(ctx *gin.Context) {
	dto := new(dtos.LoginDTO)
	if err := ctx.BindJSON(dto); err != nil {
		c.logger.Println("[ERROR][LOGIN] Error binding JSON:", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload", "Invalid request payload")
		return 
	}

	if err := c.validator.Struct(dto); err != nil {
		c.logger.Println("[ERROR][LOGIN] Validation error:", err)
		message := utils.ValidationErrorResponse(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, message, message)
		return 
	}
	token, err := c.authService.Login(dto.Email, dto.Password)
	if err != nil {
		c.logger.Println("[ERROR][LOGIN] Service error:", err)
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authentication failed", err.Error())
		return 
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Login successful", gin.H{"token": token})
}
