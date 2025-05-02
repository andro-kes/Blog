package controllers_helpers

import (
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func Connect_db(c *gin.Context) *gorm.DB{
	dbValue, ok := c.Get("DB")
	if ok == false {
		c.JSON(400, gin.H{"error": "База данных не найдена"})
		return nil
	}

	DB, ok := dbValue.(*gorm.DB)
	if ok == false {
		c.JSON(400, gin.H{"error": "Не удалось подключиться к базе данных"})
		return nil
	}
	return DB
}
	