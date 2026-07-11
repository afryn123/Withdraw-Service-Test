package infrastructure

import (
	"context"

	gormdb "github.com/afryn123/withdraw-service-test/internal/shared/gorm"
	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) Create(ctx context.Context, wallet *walletdomain.Wallet) error {
	return gormdb.DB(ctx, r.db).Create(wallet).Error
}

func (r *GormRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*walletdomain.Wallet, error) {
	return r.find(ctx, userID, false)
}

func (r *GormRepository) FindByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*walletdomain.Wallet, error) {
	return r.find(ctx, userID, true)
}

func (r *GormRepository) find(ctx context.Context, userID uuid.UUID, lock bool) (*walletdomain.Wallet, error) {
	query := gormdb.DB(ctx, r.db).Preload("User")
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var wallet walletdomain.Wallet
	if err := query.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *GormRepository) Update(ctx context.Context, wallet *walletdomain.Wallet) error {
	return gormdb.DB(ctx, r.db).
		Model(&walletdomain.Wallet{}).
		Where("id = ?", wallet.ID).
		Update("balance", wallet.Balance).Error
}

func (r *GormRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return gormdb.DB(ctx, r.db).
		Where("user_id = ?", userID).
		Delete(&walletdomain.Wallet{}).Error
}
