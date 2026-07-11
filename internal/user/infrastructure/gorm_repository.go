package infrastructure

import (
	"context"
	"errors"
	"fmt"

	gormdb "github.com/afryn123/withdraw-service-test/internal/shared/gorm"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) Create(ctx context.Context, user *userdomain.User) error {
	return gormdb.DB(ctx, r.db).Create(user).Error
}

func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	var user userdomain.User
	if err := gormdb.DB(ctx, r.db).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

func (r *GormRepository) FindByID(ctx context.Context, userID uuid.UUID) (*userdomain.User, error) {
	var user userdomain.User
	if err := gormdb.DB(ctx, r.db).First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userdomain.ErrNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *GormRepository) Update(ctx context.Context, user *userdomain.User) error {
	result := gormdb.DB(ctx, r.db).
		Model(&userdomain.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"name":       user.Name,
			"username":   user.Username,
			"email":      user.Email,
			"password":   user.Password,
			"updated_by": user.UpdatedBy,
		})
	if result.Error != nil {
		return fmt.Errorf("update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return userdomain.ErrNotFound
	}
	return nil
}

func (r *GormRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	result := gormdb.DB(ctx, r.db).Delete(&userdomain.User{}, "id = ?", userID)
	if result.Error != nil {
		return fmt.Errorf("delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return userdomain.ErrNotFound
	}
	return nil
}
