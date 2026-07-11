package infrastructure

import (
	"context"

	gormdb "github.com/afryn123/withdraw-service-test/internal/shared/gorm"
	transactiondomain "github.com/afryn123/withdraw-service-test/internal/transactions/domain"
	"gorm.io/gorm"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) Create(ctx context.Context, transaction *transactiondomain.Transaction) error {
	return gormdb.DB(ctx, r.db).Create(transaction).Error
}
