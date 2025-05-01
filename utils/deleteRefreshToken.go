package utils

import (
	"log"

	"github.com/andro-kes/Blog/models"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

func DeleteRefreshToken(DB *gorm.DB, userID uint, tokenID uuid.UUID) (error) {
	var token models.RefreshTokens
	obj := DB.Delete(&token, userID)
	log.Println("Удаление рефреш токена")
	return obj.Error
}