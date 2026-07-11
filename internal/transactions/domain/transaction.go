package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TypeWithdraw    = "withdraw"
	StatusCompleted = uint16(1)
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WalletID        uuid.UUID `gorm:"type:uuid;not null"`
	Amount          int64     `gorm:"not null"`
	Type            string    `gorm:"type:varchar(10);not null"`
	TransactionCode string    `gorm:"unique;not null"`
	ReferenceNumber string    `gorm:"unique;not null"`
	Status          uint16    `gorm:"not null"`
	Remark          *string   `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (Transaction) TableName() string { return "transaction_histories" }

type Repository interface {
	Create(ctx context.Context, transaction *Transaction) error
}

type CodeGenerator interface {
	TransactionCode() string
	ReferenceNumber(userID uuid.UUID) string
}
