package controllers

import (
	// "log"
	"time"

	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/utils"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func RefreshTokenHandler(c *gin.Context) {
	dbValue, ok := c.Get("DB")
	if ok == false {
		c.JSON(400, gin.H{"error": "База данных не найдена"})
		return
	}
	DB, ok := dbValue.(*gorm.DB)
	if ok == false {
		c.JSON(400, gin.H{"error": "Не удалось подключиться к базе данных"})
		return
	}
	
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(400, gin.H{"error": "Refresh токен не найден"})
		return
	}

	_, err = c.Cookie("token")
	if err != nil {
		c.JSON(400, gin.H{"error": "Токен не найден"})
		return
	}

	claims, err := utils.ParseRefreshToken(refreshToken)
	ok = utils.CompareTokens(DB, claims.TokenID, refreshToken)
	if ok == false {
		c.JSON(400, gin.H{"error": "Токены не совпали"})
		return
	}

	var user models.Users
	DB.Where("id = ?", claims.UserID).First(&user)
	if user.ID == 0 {
		c.JSON(400, gin.H{"error": "Пользователя не существует"})
		return
	}

	refreshToken, err = utils.UpdateRefreshToken(DB, claims.UserID, claims.TokenID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Не удалось обновить refresh токен"})
		return
	}

	tokenString, err := utils.GenerateAccessToken(user)
	if err != nil {
		c.JSON(400, gin.H{"error": "Не удалось обновить access токен"})
		return
	}

	c.SetCookie("refresh_token", refreshToken, int(time.Now().Add(7 * 24 * time.Hour).Unix()), "/", "localhost", false, true)
	c.SetCookie("token", tokenString, int(time.Now().Add(5 * time.Minute).Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Токены обновлены"})
}