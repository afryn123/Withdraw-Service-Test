package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionHistories struct {
	ID              string         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	WalletID        uuid.UUID      `gorm:"type:uuid;not null"`
	Wallet          Wallets        `gorm:"foreignKey:WalletID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount          int64          `gorm:"not null"`
	Type            string         `gorm:"type:varchar(10);not null"` // "deposit" or "withdraw"
	TransactionCode string         `gorm:"unique;not null"`
	ReferenceNumber string         `gorm:"unique;not null"`
	Status          uint16         `gorm:"not null"` // 0: pending, 1: completed, 2: failed
	Remark          *string        `gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (TransactionHistories) TableName() string {
	return "transaction_histories"
}
