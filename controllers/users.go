package controllers

import (
	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	DB, err := c.Get("DB")
	if err != false {
		c.JSON(400, gin.H{"error": "Не удалось подключиться к базе данных"})
	}

	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Ошбика при связывании с моделью Users"})
	}

	var existingUser models.Users
	DB.Where("email = ?", user.Email).First(&existingUser)

}