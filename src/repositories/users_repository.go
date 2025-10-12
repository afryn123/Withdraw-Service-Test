package repositories

import (
	"afryn123/withdraw-service/src/models"
	"fmt"

	"gorm.io/gorm"
)

type UsersRepository struct{}

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{}
}

func (r *UsersRepository) FindByEmail(db *gorm.DB, email string) (*models.Users, error) {
	var user models.Users
	if err := db.
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (r *UsersRepository) Create(db *gorm.DB, user *models.Users) error {
	return db.Create(user).Error
}