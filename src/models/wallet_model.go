package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallets struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid;unique;not null"`
	Users     Users          `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Balance   int64          `gorm:"not null;default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Wallets) TableName() string {
	return "wallets"
}
