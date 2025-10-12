package repositories

import (
	"afryn123/withdraw-service/src/models"

	"gorm.io/gorm"
)

type TransactionHistoryRepository struct{}

func NewTransactionHistoryRepository() *TransactionHistoryRepository {
	return &TransactionHistoryRepository{}
}

func (r *TransactionHistoryRepository) Create(db *gorm.DB, transaction *models.TransactionHistories) error {
	return db.Create(&transaction).Error
}
