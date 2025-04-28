package utils

import (
	"github.com/andro-kes/Blog/models"

	"github.com/google/uuid"

	"gorm.io/gorm"

	"log"
)

func CompareTokens(DB *gorm.DB, tokenID uuid.UUID, token string) bool {
	var existingToken models.RefreshTokens
	DB.Where("token_id = ?", tokenID).First(&existingToken)
	
	if existingToken.Token == "" {
		log.Println("Токен не найден")
		return false
	}

	if existingToken.Token == token {
		return true
	}
	return false
}