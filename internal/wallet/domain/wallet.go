package domain

import (
	"context"
	"time"

	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID       `gorm:"type:uuid;unique;not null"`
	User      userdomain.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Balance   int64           `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Wallet) TableName() string { return "wallets" }

func (w *Wallet) Withdraw(amount int64) error {
	if amount <= 0 {
		return ErrInvalidWithdrawAmount
	}
	if w.Balance < amount {
		return ErrInsufficientBalance
	}
	w.Balance -= amount
	return nil
}

type Repository interface {
	Create(ctx context.Context, wallet *Wallet) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	FindByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	Update(ctx context.Context, wallet *Wallet) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
