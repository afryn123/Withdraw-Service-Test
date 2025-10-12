package controllers

import (
	"afryn123/withdraw-service/src/dtos"
	"afryn123/withdraw-service/src/services"
	"afryn123/withdraw-service/src/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TransactionContoller struct {
	transactionServive *services.TransactionService
	logger             *log.Logger
	validator          *validator.Validate
}

func NewTransactionContoller(transactionServive *services.TransactionService) *TransactionContoller {
	return &TransactionContoller{
		transactionServive: transactionServive,
		logger:             utils.GetLogger("txn"),
		validator:          validator.New(),
	}
}

func (c *TransactionContoller) Withdraw(ctx *gin.Context) {
	dto := new(dtos.WithdrawRequestDTO)
	if err := ctx.BindJSON(dto); err != nil {
		c.logger.Println("[ERROR][WITHDRAW] Error binding JSON:", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload", "Invalid request payload")
		return
	}
	if err := c.validator.Struct(dto); err != nil {
		message := utils.ValidationErrorResponse(err)
		c.logger.Println("[ERROR][WITHDRAW] Validation error:", message)
		utils.ErrorResponse(ctx, http.StatusBadRequest, message, nil)
		return
	}
	userId, exists := ctx.Get("userID")
	if !exists {
		c.logger.Println("[ERROR][WITHDRAW] User ID not found in context")
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "User not authenticated")
		return
	}

	userIDStr, ok := userId.(string)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Invalid user ID format", nil)
		return
	}

	userIDParse, err := uuid.Parse(userIDStr)
	if err != nil {
		c.logger.Println("[ERROR][WITHDRAW] Invalid UUID format:", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID format", "Invalid user ID format")
		return
	}

	result, err := c.transactionServive.Withdraw(userIDParse, dto.Amount, dto.Remark)
	if err != nil {
		c.logger.Println("[ERROR][WITHDRAW] Service error:", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to process withdraw", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Withdraw processed successfully", result)
}
