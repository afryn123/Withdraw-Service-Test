package config

import (
	"afryn123/withdraw-service/src/dtos"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(env *dtos.EnvirontmentVariable) {
	// development
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s sslmode=%s TimeZone=Asia/Jakarta",
		env.DBHost,
		env.DBUser,
		env.DBName,
		env.DBPort,
		env.DBPass,
		env.DBSSL,
	)

	log.Println("DSN:", dsn) // Debug DSN

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // slow SQL Threshold
			LogLevel:      logger.Info, // log level
			Colorful:      true,        // colorize
		},
	)

	var err error
	if (env.ApiEnv == "production") {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger})
	}

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Check connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic DB: %v", err)
	}

	// Set connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected and migration successful")
}
