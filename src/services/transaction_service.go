package services

import (
	"afryn123/withdraw-service/src/dtos"
	"afryn123/withdraw-service/src/models"
	"afryn123/withdraw-service/src/repositories"
	"afryn123/withdraw-service/src/utils"
	"errors"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	db              *gorm.DB
	walletRepo      *repositories.WalletRepository
	transactionRepo *repositories.TransactionHistoryRepository
	logger          *log.Logger
}

func NewTransactionService(
	db *gorm.DB,
	walletRepo *repositories.WalletRepository,
	transactionRepo *repositories.TransactionHistoryRepository,
) *TransactionService {
	return &TransactionService{
		db:              db,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		logger:          utils.GetLogger("txn"),
	}
}

func (s *TransactionService) Withdraw(userId uuid.UUID, amount int64, remark *string) (dtos.WithdrawResponseDTO, error) {
	var err error
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			s.logger.Println("[ERROR][TRANSACTION WITHDRAW] Transaction panic, rolling back:", r)
			tx.Rollback()
		}
	}()
	

	err = s.walletRepo.LockRowForUpdate(tx, &models.Wallets{UserID: userId})
	if err != nil {
		s.logger.Println("[ERROR][TRANSACTION WITHDRAW] Lock row error:", err)
		tx.Rollback()
		return dtos.WithdrawResponseDTO{}, err
	}

	wallet, err := s.walletRepo.FindByUserId(tx, userId)
	if err != nil {
		return dtos.WithdrawResponseDTO{}, err
	}

	if wallet.Balance < amount {
		return dtos.WithdrawResponseDTO{}, errors.New("insufficient balance")
	}

	wallet.Balance -= amount

	err = s.walletRepo.UpdateWithLocking(tx, &wallet)
	if err != nil {
		s.logger.Println("[ERROR][TRANSACTION WITHDRAW] Update wallet error:", err)
		tx.Rollback()
		return dtos.WithdrawResponseDTO{}, err
	}
	
	transaction := &models.TransactionHistories{
		WalletID:      wallet.ID,
		Amount:      amount,
		Type:        "withdraw",
		TransactionCode: utils.GenerateTransactionCode(),
		ReferenceNumber: utils.GenerateReferenceNumber(wallet.Users.ID.String()),
		Status:      1, // completed
		Remark:      remark,
	}
	err = s.transactionRepo.Create(tx, transaction)
	if err != nil {
		s.logger.Println("[ERROR][TRANSACTION WITHDRAW] Create transaction error:", err)
		tx.Rollback()
		return dtos.WithdrawResponseDTO{}, err
	}
	
	if err := tx.Commit().Error; err != nil {
		s.logger.Println("[ERROR][TRANSACTION WITHDRAW] Commit transaction error:", err)
		return dtos.WithdrawResponseDTO{}, err
	}
	
	response := dtos.WithdrawResponseDTO{
		TransactionID:   transaction.ID,
		ReferenceNumber: transaction.ReferenceNumber,
		Transaction:  dtos.SubTransactionWithdrawResponseDto{
			Amount:   transaction.Amount,
			Type:    transaction.Type,
			BalanceNow: wallet.Balance,
			Remark: remark,
		},
		User: dtos.SubUserResponseDto{
			UserID: wallet.Users.ID.String(),
			Name:   wallet.Users.Name,
		},
	}
	return response, nil

}
