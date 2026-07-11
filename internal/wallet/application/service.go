package application

import (
	"context"

	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
	"github.com/google/uuid"
)

type Balance struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Balance  int64     `json:"balance"`
	User     User      `json:"user"`
}

type User struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

type Service struct{ wallets walletdomain.Repository }

func NewService(wallets walletdomain.Repository) *Service { return &Service{wallets: wallets} }

func (s *Service) FindBalanceByUserID(ctx context.Context, userID uuid.UUID) (Balance, error) {
	wallet, err := s.wallets.FindByUserID(ctx, userID)
	if err != nil {
		return Balance{}, err
	}
	return Balance{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
		User:     User{UserID: wallet.User.ID, Name: wallet.User.Name},
	}, nil
}
