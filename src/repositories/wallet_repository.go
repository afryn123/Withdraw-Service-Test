package repositories

import (
	"afryn123/withdraw-service/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct{}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{}
}

func (r *WalletRepository) FindByUserId(db *gorm.DB, userId uuid.UUID) (models.Wallets, error) {
	var data models.Wallets
	if err := db.
		Preload("Users").
		Where("user_id = ?", userId).
		First(&data).Error; err != nil {
		return models.Wallets{}, err
	}

	return data, nil
}

// Create new wallet
func (r *WalletRepository) Create(db *gorm.DB, wallet *models.Wallets) error {
	return db.Create(&wallet).Error
}

// Lock row for update
func (r *WalletRepository) LockRowForUpdate(db *gorm.DB, wallet *models.Wallets) error	{
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", wallet.UserID).
		First(&wallet).Error; err != nil {
		db.Rollback()
		return err
	}
	return nil
}

// Update wallet with locking
func (r *WalletRepository) Update(db *gorm.DB, wallet *models.Wallets) error {
	
	return db.Save(&wallet).Error
}