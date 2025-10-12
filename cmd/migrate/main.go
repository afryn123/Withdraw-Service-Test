package main

import (
	"afryn123/withdraw-service/src/config"
	"afryn123/withdraw-service/src/utils"
	"log"
)

func main() {
	// Load environment variables
	env := utils.Environtment()

	// Connect to DB
	config.ConnectDatabase(env)

	if config.DB == nil {
		log.Fatal("Database initialization failed")
	}

	// Migrate tables
	utils.MigrateTable()
}