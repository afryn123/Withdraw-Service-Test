package main

import (
	"os"

	"github.com/afryn123/withdraw-service-test/internal/shared/config"
	gormdb "github.com/afryn123/withdraw-service-test/internal/shared/gorm"
	sharedlogger "github.com/afryn123/withdraw-service-test/internal/shared/logger"
	transactiondomain "github.com/afryn123/withdraw-service-test/internal/transactions/domain"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
)

func main() {
	cfg := config.Load()
	logger := sharedlogger.New(cfg.LogLevel).With("service", "withdraw-migrate", "environment", cfg.Environment)
	db, err := gormdb.Open(cfg.DSN(), logger)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	if err := db.AutoMigrate(&userdomain.User{}, &walletdomain.Wallet{}, &transactiondomain.Transaction{}); err != nil {
		logger.Error("database migration failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database migration completed")
}
