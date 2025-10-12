package services

import (
	"afryn123/withdraw-service/src/repositories"
	"afryn123/withdraw-service/src/utils"
	"fmt"

	"gorm.io/gorm"
)

type AuthService struct {
	db       *gorm.DB
	userRepo *repositories.UsersRepository
}

func NewAuthService(db *gorm.DB, userRepo *repositories.UsersRepository) *AuthService {
	return &AuthService{
		db:       db,
		userRepo: userRepo,
	}
}

// login
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(s.db, email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", fmt.Errorf("invalid password")
	}

	token, _ := utils.GenerateJWT(user.ID.String())
	return token, nil
}
