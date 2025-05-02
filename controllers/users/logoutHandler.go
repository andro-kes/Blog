package users_controllers

import (
	"log"

	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/utils"
	"github.com/andro-kes/Blog/controllers/helpers"
	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	DB := controllers_helpers.Connect_db(c)
	if DB == nil {
		return
	}

	token, err := c.Cookie("refresh_token")
	if err != nil {
		log.Println("LogoutHandler: Refresh токен не найден в cookie")
		c.JSON(400, gin.H{"error": "Refresh токен не найден"})
	}
	
	claims, err := utils.ParseRefreshToken(token)
	if err != nil {
		log.Println("LogoutHandler: Claims не определены")
		c.JSON(400, gin.H{"error": "Права пользователя не определены"})
	}

	var RefreshToken models.RefreshTokens
	obj := DB.Delete(&RefreshToken, claims.UserID)
	if obj.Error != nil {
		log.Println("LogoutHandler: Refresh токен не найден в базе")
		c.JSON(400, gin.H{"error": "Ошибка при получении refresh еокена из базы данных"})
	}

	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Пользователь вышел из системы"})
}