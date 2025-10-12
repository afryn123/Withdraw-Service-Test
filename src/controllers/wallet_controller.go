package controllers

import (
	"afryn123/withdraw-service/src/services"
	"afryn123/withdraw-service/src/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletContoller struct {
	walletService *services.WalletService
	logger        *log.Logger
}

func NewWalletContoller(walletService *services.WalletService) *WalletContoller {
	return &WalletContoller{
		walletService: walletService,
		logger:        utils.GetLogger("system"),
	}
}

func (c *WalletContoller) GetBalanceByUserId(ctx *gin.Context) {
	strUserID := ctx.Param("user_id")
	userId, err := uuid.Parse(strUserID)
	if err != nil {
		c.logger.Println("[ERROR][GET BALANCE] Invalid UUID format:", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID format", "Invalid user ID format")
		return
	}
	
	balanceResponse, err := c.walletService.FindBalanceByUserId(userId)
	if err != nil {
		c.logger.Println("[ERROR][GET BALANCE] Service error:", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve balance", "Failed to retrieve balance")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Balance retrieved successfully", balanceResponse)
}
