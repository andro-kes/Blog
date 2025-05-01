package controllers

import (
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/utils"
	
	"github.com/gin-gonic/gin"
)

func SignupHandler(c *gin.Context) {
	DB := connect_db(c)
	if DB == nil {
		return
	}

	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Ошибка при связывании с моделью Users"})
		return
	}

	var existingUser models.Users
	DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		c.JSON(400, gin.H{"error": "Пользователь с такими данными уже существует"})
		return
	}

	hashPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": "Не удалось получить хэш пароля"})
		return
	}
	user.Password = string(hashPassword)

	DB.Create(&user)
	c.JSON(201, gin.H{"message": "Пользователь успешно зарегистрирован"})
}