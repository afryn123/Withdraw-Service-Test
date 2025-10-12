package routes

import (
	"afryn123/withdraw-service/src/controllers"
	"afryn123/withdraw-service/src/middlewares"
	"afryn123/withdraw-service/src/repositories"
	"afryn123/withdraw-service/src/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Routes struct {
	db                    *gorm.DB
	route                 *gin.Engine
	userController        *controllers.UsersContoller
	walletController      *controllers.WalletContoller
	transactionController *controllers.TransactionContoller
	authController        *controllers.AuthContoller
}

func NewRoutes(db *gorm.DB, route *gin.Engine) *Routes {
	// repositories
	userRepo := repositories.NewUsersRepository()
	walletRepo := repositories.NewWalletRepository()
	transactionHistoriesRepo := repositories.NewTransactionHistoryRepository()

	// services
	userService := services.NewUsersService(db, userRepo, walletRepo)
	walletService := services.NewWalletService(db, walletRepo)
	transactionService := services.NewTransactionService(db, walletRepo, transactionHistoriesRepo)
	authService := services.NewAuthService(db, userRepo)

	// controllers
	walletController := controllers.NewWalletContoller(walletService)
	userController := controllers.NewUsersContoller(userService)
	transactionController := controllers.NewTransactionContoller(transactionService)
	authController := controllers.NewAuthController(authService)

	return &Routes{
		db:                    db,
		route:                 route,
		userController:        userController,
		walletController:      walletController,
		transactionController: transactionController,
		authController:        authController,
	}
}

func (r *Routes) AuthRouter() {

	api := r.route.Group("/api/auth")

	// Auth routes
	api.POST("/login", r.authController.Login)
}

func (r *Routes) UserRouter() {

	api := r.route.Group("/api/users", middlewares.AuthProtected())

	// User routes
	api.POST("/create", r.userController.Create)
}

func (r *Routes) WalletRouter() {
	api := r.route.Group("/api/wallet", middlewares.AuthProtected())

	// Wallet routes
	api.GET(":user_id/balance", r.walletController.GetBalanceByUserId)
}
func (r *Routes) TransactionRouter() {
	api := r.route.Group("/api/transaction", middlewares.AuthProtected())

	// Wallet routes
	api.POST("/withdraw", r.transactionController.Withdraw)
}
