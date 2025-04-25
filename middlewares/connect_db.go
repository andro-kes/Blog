package middlewares

import (
	"fmt"
	"os"

	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "Файл .env не открывается"})
		}
		dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("HOST"),
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("PORT"),
		os.Getenv("SSLMODE"),
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "Ошибка при открытии базы данных"})
		}
		if err := db.AutoMigrate(&models.Users{}); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "Ошибка миграции"})
		}

		c.Set("DB", db)
		c.Next()
	}
}