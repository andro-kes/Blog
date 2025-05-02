package utils

import (
	"github.com/andro-kes/Blog/config"
	"github.com/andro-kes/Blog/models"

	"github.com/google/uuid"
	"github.com/dgrijalva/jwt-go"

	"gorm.io/gorm"

	"time"
	"strconv"
	"log"
)

func GenerateRefreshToken(DB *gorm.DB, userID uint) (string, error) {
	log.Println("Генерация рефреш токена")
	expititionTime := time.Now().Add(7 * 24 * time.Hour)

	TokenID := uuid.New()

	refreshClaims := models.RefreshTokensClaims{
		UserID: userID,
		TokenID: TokenID,
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(int(userID)),
			ExpiresAt: expititionTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	tokenString, err := token.SignedString([]byte(config.SECRET_KEY))
	
	var existingToken models.RefreshTokens
	DB.Where("user_id = ?", userID).First(&existingToken)
	if existingToken.Token != "" {
		log.Println("Обновление реферш токен")
		DB.Model(&existingToken).Update("token", tokenString)
		return tokenString, err
	}

	refreshToken := &models.RefreshTokens{
		UserID: userID,
		TokenID: TokenID,
		Token: tokenString,
	}

	log.Println("Добавления рефреш токена в базу", refreshToken.TokenID)
	obj := DB.Create(&refreshToken)
	if obj.Error != nil {
		log.Println("Ошибка при добавлении рефреш токена в базу")
		return "", obj.Error
	}

	return tokenString, err
}