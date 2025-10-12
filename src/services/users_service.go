package services

import (
	"afryn123/withdraw-service/src/dtos"
	"afryn123/withdraw-service/src/models"
	"afryn123/withdraw-service/src/repositories"
	"afryn123/withdraw-service/src/utils"
	"log"

	"gorm.io/gorm"
)

type UsersService struct {
	db        *gorm.DB
	usersRepo *repositories.UsersRepository
	walletRepo *repositories.WalletRepository
	logger    *log.Logger
}

func NewUsersService(
	db *gorm.DB,
	usersRepo *repositories.UsersRepository,
	walletRepo *repositories.WalletRepository,
) *UsersService {
	return &UsersService{
		db:        db,
		usersRepo: usersRepo,
		walletRepo: walletRepo,
		logger:    utils.GetLogger("system"),
	}
}

func (s *UsersService) Create(dto *dtos.UserCreateDTO) error {

	hashedPassword, err := utils.HashPassword(dto.Password)

	if err != nil {
		log.Println("Failed to hash password:", err)
		return err
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			s.logger.Println("Transaction panic, rolling back:", r)
			tx.Rollback()
		}
	}()
	

	user := &models.Users{
		Name:      dto.Name,
		Username:  dto.Username,
		Email:     dto.Email,
		Password:  hashedPassword,
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.CreatedBy,
	}
	if err := s.usersRepo.Create(tx, user); err != nil { 
		tx.Rollback()
		return err
	}

	wallet := &models.Wallets{
		UserID:    user.ID,
		Balance:   0,
	} 

	if err := s.walletRepo.Create(tx, wallet); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
