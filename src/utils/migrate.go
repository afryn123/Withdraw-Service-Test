package utils

import (
	"afryn123/withdraw-service/src/config"
	"afryn123/withdraw-service/src/models"
	"log"
)

func MigrateTable() {
	// Migrate
	if err := config.DB.AutoMigrate(
		&models.Users{},
		&models.Wallets{},
		&models.TransactionHistories{},

	); err != nil {
		log.Fatalf(" Error migrating database: %v", err)
	}
}	