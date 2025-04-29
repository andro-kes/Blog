package middlewares

import (
	"log"
	"time"

	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"
	"github.com/andro-kes/Blog/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var TestDB *gorm.DB

func init() { 
	config.LoadConfig()
	
	DB = openDB(config.DSN)
	TestDB = openDB(config.TESTDSN)
}

func openDB(DSN string) *gorm.DB {
	var err error
	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err) 
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Ошибка при получении *sql.DB: %v", err)
		return nil
	}

    sqlDB.SetMaxIdleConns(10)  
    sqlDB.SetMaxOpenConns(100)   
    sqlDB.SetConnMaxLifetime(time.Hour) 

	if err := DB.AutoMigrate(&models.Users{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
		return nil 
	}
	if err := DB.AutoMigrate(&models.RefreshTokens{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
		return nil
	}
	log.Println("Успешное подключение к базе данных и миграция выполнены")
	return DB
}

func DBMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if DB == nil || TestDB == nil {
			log.Println("DBMiddleWare: База данных не инициализирована")
			c.AbortWithStatusJSON(500, gin.H{"error": "База данных не инициализирована"})
			return
		}
		log.Println(gin.Mode())
		if gin.Mode() == "test" {
			c.Set("DB", TestDB)
		} else {
			c.Set("DB", DB)
		}
		c.Next()
	}
}