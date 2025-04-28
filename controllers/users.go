package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/andro-kes/Blog/config"
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"

	"gorm.io/gorm"
)

func init() {
	config.LoadConfig()
	oauth2Config = &oauth2.Config{
		ClientID: config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		RedirectURL: config.REDIRECT_URL,
		Scopes: []string{"login:email"},
		Endpoint: yandex.Endpoint,
	}
}

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

	expititionTime := time.Now().Add(5 * time.Minute)
	claims := models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject: existingUser.Email,
			ExpiresAt: expititionTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SECRET_KEY))
	if err != nil {
		c.JSON(400, gin.H{"error": "Не удалось создать токен"})
		return
	}

	c.SetCookie("token", tokenString, int(expititionTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Пользователь авторизован"})
}

func SignupHandler(c *gin.Context) {
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
	c.JSON(300, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

func LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Пользователь вышел из системы"})
}


type YandexUser struct {
	DefaultEmail string `json:"default_email"`
	Emails []string `json:"emails"`
}

var (
	oauth2Config *oauth2.Config
	oauthStateString = "random-string-for-state" // TODO Защита от CSRF
)

func AuthYandexRedirectHandler(c *gin.Context) {
	url := oauth2Config.AuthCodeURL(oauthStateString)
	c.Redirect(307, url)
}

func LoginYandexHandler(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		log.Printf("Неверный статус, ожидалось: '%s', получено: '%s'\n", oauthStateString, state)
		c.AbortWithStatus(401)
		return
	}

	code := c.Query("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("oauthConf.Exchange() не сработал из-за: '%s'\n", err)
		c.AbortWithStatus(400)
		return
	}

	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://login.yandex.ru/info?format=json")
	if err != nil {
		log.Printf("Не удалось получить информацию о пользователе: %v", err)
		c.AbortWithStatus(500)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Не удалось прочитать тело ответа: %v", err)
		c.AbortWithStatus(500)
		return
	}

	var yandexUser YandexUser
	err = json.Unmarshal(body, &yandexUser)
	if err != nil {
		log.Printf("Не удалось разобрать JSON: %v", err)
		c.AbortWithStatus(500)
		return
	}

	log.Printf("Информация о пользователе (email): %+v\n", yandexUser)

	var existingUser models.Users
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
	
	DB.Where("email = ?", yandexUser.DefaultEmail).First(&existingUser)
	if existingUser.ID != 0 {
		log.Println("Вход в систему через Яндекс")
	} else {
		DB.Create(&models.Users{
			Email: yandexUser.DefaultEmail,
			Role: "user",
			IsOauth: true,
		})
	}
	// TODO вынести в функцию создания токена
	expititionTime := time.Now().Add(5 * time.Minute)
	claims := models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject: existingUser.Email,
			ExpiresAt: expititionTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(config.SECRET_KEY))
	if err != nil {
		c.JSON(400, gin.H{"error": "Не удалось создать токен"})
		return
	}

	c.SetCookie("token", tokenString, int(expititionTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"message": "Успешная авторизация через Яндекс",
		"email":   yandexUser.DefaultEmail,
	})
}