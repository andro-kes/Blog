package controllers

import (
	"time"

	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/utils"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func LoginHandler(c *gin.Context) {
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

	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Ошибка при связывании с моделью Users"})
		return
	}

	var existingUser models.Users
	if err := DB.Where("email = ?", user.Email).First(&existingUser); err.Error != nil {
		c.JSON(400, gin.H{"error": "Ошибка при обращении к базе"})
		return
	}

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "Такого пользователя не существует"})
		return
	}

	if err := utils.CompareHashPasswords(user.Password, existingUser.Password); err != nil {
		c.JSON(400, gin.H{"error": "Неверный пароль"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(DB, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Не удалось сгенерировать refresh токен"})
		return
	}

	tokenString, err := utils.GenerateAccessToken(existingUser)
	if err != nil {
		c.JSON(400, gin.H{"error": "Не удалось сгенерировать access токен"})
		return
	}
	expititionTime := time.Now().Add(5 * time.Minute)

	c.SetCookie("refresh_token", refreshToken, int(time.Now().Add(7 * 24 * time.Hour).Unix()), "/", "localhost", false, true)
	c.SetCookie("token", tokenString, int(expititionTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"message": "Пользователь авторизован",
		"email": user.Email,
	})
}