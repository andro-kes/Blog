package config

import (
	"log"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var SECRET_KEY string
var DSN string

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err) 
		return
	}

	DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("HOST"),
		"postgres", // TODO: Перенести в .env
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("PORT"),
		os.Getenv("SSLMODE"),
	)

	SECRET_KEY = os.Getenv("SECRET_KEY")
	log.Println("Конфигурация установлена")
}