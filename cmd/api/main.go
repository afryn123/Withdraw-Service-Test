package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authapp "github.com/afryn123/withdraw-service-test/internal/auth/application"
	authdelivery "github.com/afryn123/withdraw-service-test/internal/auth/delivery"
	authinfra "github.com/afryn123/withdraw-service-test/internal/auth/infrastructure"
	"github.com/afryn123/withdraw-service-test/internal/shared/config"
	gormdb "github.com/afryn123/withdraw-service-test/internal/shared/gorm"
	sharedlogger "github.com/afryn123/withdraw-service-test/internal/shared/logger"
	sharedmiddleware "github.com/afryn123/withdraw-service-test/internal/shared/middleware"
	sharedopenapi "github.com/afryn123/withdraw-service-test/internal/shared/openapi"
	transactionapp "github.com/afryn123/withdraw-service-test/internal/transactions/application"
	transactiondelivery "github.com/afryn123/withdraw-service-test/internal/transactions/delivery"
	transactioninfra "github.com/afryn123/withdraw-service-test/internal/transactions/infrastructure"
	userapp "github.com/afryn123/withdraw-service-test/internal/user/application"
	userdelivery "github.com/afryn123/withdraw-service-test/internal/user/delivery"
	userinfra "github.com/afryn123/withdraw-service-test/internal/user/infrastructure"
	walletapp "github.com/afryn123/withdraw-service-test/internal/wallet/application"
	walletdelivery "github.com/afryn123/withdraw-service-test/internal/wallet/delivery"
	walletinfra "github.com/afryn123/withdraw-service-test/internal/wallet/infrastructure"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	logger := sharedlogger.New(cfg.LogLevel).With("service", "withdraw-api", "environment", cfg.Environment)
	if err := run(cfg, logger); err != nil {
		logger.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run(cfg config.Config, logger *slog.Logger) error {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	db, err := gormdb.Open(cfg.DSN(), logger)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database connection: %w", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}()

	users := userinfra.NewGormRepository(db)
	wallets := walletinfra.NewGormRepository(db)
	transactions := transactioninfra.NewGormRepository(db)
	txManager := gormdb.NewTransactionManager(db)
	passwords := userinfra.NewBcryptPassword(bcrypt.DefaultCost)
	tokens := authinfra.NewJWTManager(cfg.JWTSecret, 24*time.Hour)
	codes := transactioninfra.NewCodeGenerator()

	userHandler := userdelivery.NewHandler(userapp.NewService(users, wallets, passwords, txManager))
	walletHandler := walletdelivery.NewHandler(walletapp.NewService(wallets))
	transactionHandler := transactiondelivery.NewHandler(transactionapp.NewService(wallets, transactions, txManager, codes))
	authHandler := authdelivery.NewHandler(authapp.NewService(users, passwords, tokens))

	router := gin.New()
	router.Use(sharedmiddleware.RequestLogger(logger), sharedmiddleware.CORS(), sharedmiddleware.Recovery(logger))
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "OK"}) })
	sharedopenapi.RegisterRoutes(router, sharedopenapi.NewHandler())
	userdelivery.RegisterPublicRoutes(router.Group("/api/users"), userHandler)
	authdelivery.RegisterRoutes(router.Group("/api/auth"), authHandler)
	protected := authdelivery.Middleware(tokens)
	userdelivery.RegisterProtectedRoutes(router.Group("/api/users", protected), userHandler)
	walletdelivery.RegisterRoutes(router.Group("/api/wallet", protected), walletHandler)
	transactiondelivery.RegisterRoutes(router.Group("/api/transaction", protected), transactionHandler)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()
	logger.Info("HTTP server started", "address", server.Addr)

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownSignals)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("run server: %w", err)
		}
		return nil
	case receivedSignal := <-shutdownSignals:
		logger.Info("shutdown signal received", "signal", receivedSignal.String())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown timed out", "error", err)
		if closeErr := server.Close(); closeErr != nil {
			logger.Error("failed to force close server", "error", closeErr)
		}
	}

	if err := <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server stopped: %w", err)
	}
	logger.Info("HTTP server stopped")
	return nil
}
