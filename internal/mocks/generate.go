//go:generate go run go.uber.org/mock/mockgen -source=../user/domain/user.go -destination=user.go -package=mocks -mock_names=Repository=MockUserRepository,PasswordHasher=MockPasswordHasher
//go:generate go run go.uber.org/mock/mockgen -source=../wallet/domain/wallet.go -destination=wallet.go -package=mocks -mock_names=Repository=MockWalletRepository
//go:generate go run go.uber.org/mock/mockgen -source=../transactions/domain/transaction.go -destination=transaction.go -package=mocks -mock_names=Repository=MockTransactionRepository,CodeGenerator=MockCodeGenerator
//go:generate go run go.uber.org/mock/mockgen -source=../shared/transaction/manager.go -destination=shared.go -package=mocks -mock_names=Manager=MockTransactionManager
//go:generate go run go.uber.org/mock/mockgen -source=../auth/domain/ports.go -destination=auth.go -package=mocks -mock_names=PasswordVerifier=MockPasswordVerifier,TokenManager=MockTokenManager
//go:generate go run go.uber.org/mock/mockgen -source=../user/delivery/http_handler.go -destination=user_delivery.go -package=mocks -mock_names=Service=MockUserDeliveryService
//go:generate go run go.uber.org/mock/mockgen -source=../auth/delivery/http_handler.go -destination=auth_delivery.go -package=mocks -mock_names=Service=MockAuthDeliveryService
//go:generate go run go.uber.org/mock/mockgen -source=../wallet/delivery/http_handler.go -destination=wallet_delivery.go -package=mocks -mock_names=Service=MockWalletDeliveryService
//go:generate go run go.uber.org/mock/mockgen -source=../transactions/delivery/http_handler.go -destination=transaction_delivery.go -package=mocks -mock_names=Service=MockTransactionDeliveryService

package mocks
