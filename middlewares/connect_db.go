package middlewares

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB 

func init() { 
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err) 
		return
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("HOST"),
		"postgres", // TODO: Перенести в .env
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("PORT"),
		os.Getenv("SSLMODE"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err) 
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Ошибка при получении *sql.DB: %v", err)
		return
	}

    sqlDB.SetMaxIdleConns(10)  
    sqlDB.SetMaxOpenConns(100)   
    sqlDB.SetConnMaxLifetime(time.Hour) 

	if err := DB.AutoMigrate(&models.Users{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
		return
	}
	log.Println("Успешное подключение к базе данных и миграция выполнены")
}

func DBMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if DB == nil {
			log.Println("DBMiddleWare: База данных не инициализирована")
			c.AbortWithStatusJSON(500, gin.H{"error": "База данных не инициализирована"})
			return
		}
		c.Set("DB", DB)
		c.Next()
	}
}