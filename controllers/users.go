package controllers

import (
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/utils"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func connectDB(c *gin.Context) *gorm.DB {
	dbValue, ok := c.Get("DB")
	if ok != false {
		c.JSON(400, gin.H{"error": "База данных не найдена"})
		return nil
	}

	DB, ok := dbValue.(*gorm.DB)
	if ok != false {
		c.JSON(400, gin.H{"error": "Не удалось подключиться к базе данных"})
		return nil
	}
	return DB
}

func LoginHandler(c *gin.Context) {
	DB := connectDB(c)
	if DB == nil {
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

	expititionTime := time.Now().Add(5 * time.Minute)
	claims := models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject: existingUser.Email,
			ExpiresAt: expititionTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SigningString()
	if err != nil {
		c.JSON(400, gin.H{"error": "Не удалось создать токен"})
		return
	}

	c.SetCookie("token", tokenString, int(expititionTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Пользователь авторизован"})
}

func SignupHandler(c *gin.Context) {
	DB := connectDB(c)
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

	hashPassword, err := utils.GenerateHashPassword(existingUser.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": "Не удалось получить хэш пароля"})
		return
	}
	existingUser.Password = string(hashPassword)

	DB.Create(&existingUser)
	c.JSON(300, gin.H{"message": "Пользователь успешно зарегистрирован"})
}