package controllers

import (
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/utils"

	"github.com/gin-gonic/gin"
	
	"context"
	"encoding/json"
	"io"
	"log"
	"time"
)

type YandexUser struct {
	DefaultEmail string `json:"default_email"`
	Emails []string `json:"emails"`
}

var (
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

	var existingUser models.Users
	DB := connect_db(c)
	if DB == nil {
		return
	}
	
	DB.Where("email = ?", yandexUser.DefaultEmail).First(&existingUser)
	if existingUser.ID != 0 {
		log.Println("Вход в систему через Яндекс")
	} else {
		existingUser.Email = yandexUser.DefaultEmail
		existingUser.Role = "user"
		existingUser.IsOauth = true
		DB.Create(&existingUser)
	}

	refreshToken, err := utils.GenerateRefreshToken(DB, existingUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	tokenString, err := utils.GenerateAccessToken(existingUser)
	if err != nil {
		c.JSON(400, gin.H{"error": "Ошибка при создании токена"})
		return
	}
	expititionTime := time.Now().Add(5 * time.Minute)

	c.SetCookie("refresh_token", refreshToken, int(time.Now().Add(7 * 24 * time.Hour).Unix()), "/", "localhost", false, true)
	c.SetCookie("token", tokenString, int(expititionTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"message": "Успешная авторизация через Яндекс",
		"email":   yandexUser.DefaultEmail,
	})
}