package utils

import (
	"afryn123/withdraw-service/src/dtos"
	"log"
	"os"

	"github.com/joho/godotenv"
)


func Environtment() (ev *dtos.EnvirontmentVariable) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	ev = &dtos.EnvirontmentVariable{}
	ev.ApiEnv = os.Getenv("API_ENV")
	ev.DBName = os.Getenv("DB_NAME")
	ev.DBUser = os.Getenv("DB_USER")
	ev.DBPort = os.Getenv("DB_PORT")
	ev.DBHost = os.Getenv("DB_HOST")
	ev.DBPass = os.Getenv("DB_PASS")
	ev.DBSSL = os.Getenv("DB_SSL_MODE")
	ev.AppHost = os.Getenv("APP_HOST")
	ev.AppPort = os.Getenv("APP_PORT")
	ev.JwtSecret = os.Getenv("JWT_SECRET")

	return ev
}
