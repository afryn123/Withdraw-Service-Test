package services

import (
	"afryn123/withdraw-service/src/dtos"
	"afryn123/withdraw-service/src/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService struct {
	db  *gorm.DB
	walletRepo *repositories.WalletRepository
}

func NewWalletService(db *gorm.DB, walletRepo *repositories.WalletRepository) *WalletService {
	return &WalletService{
		db: db,
		walletRepo: walletRepo,
	}
}

func (s *WalletService) FindBalanceByUserId(userID uuid.UUID) (dtos.BalanceResponseDto,error) {
	wallet, err := s.walletRepo.FindByUserId(s.db, userID)
	if err != nil {
		return dtos.BalanceResponseDto{}, err
	}

	return dtos.BalanceResponseDto{
		WalletID: wallet.ID,
		Balance: wallet.Balance,
		User: dtos.SubUserBalanceResponseDto{
			UserID: wallet.Users.ID,
			Name: wallet.Users.Name,
		},
	}, nil
}

