package main

import (
	"afryn123/withdraw-service/src/config"
	"afryn123/withdraw-service/src/middlewares"
	"afryn123/withdraw-service/src/routes"
	"afryn123/withdraw-service/src/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// environtment
	env := utils.Environtment()

	// logger
	utils.InitLoggers()
	logger := utils.GetLogger("system")
	if env.ApiEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		logger.Println("[INFO][ENV] System Started on Production")
	} else {
		logger.Println("[INFO][ENV] System started on Development")
	}

	// Database
	config.ConnectDatabase(env)
	if config.DB == nil {
		logger.Fatal("[ERROR][DB] Database initialization failed")
	}

	// Gin
	r := gin.Default()

	// CORS
	r.Use(middlewares.Cors())

	// res 500 when panic
	r.Use(middlewares.CustomRecoverPanic())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// Routes
	appRoutes := routes.NewRoutes(config.DB, r)
	appRoutes.UserRouter()
	appRoutes.WalletRouter()
	appRoutes.TransactionRouter()
	appRoutes.AuthRouter()

	port := env.AppPort
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
