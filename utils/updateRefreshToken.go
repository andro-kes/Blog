package utils

import (
	"github.com/andro-kes/Blog/models"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

func UpdateRefreshToken(DB *gorm.DB, userID uint, tokenID uuid.UUID) (string, error) {
	var token models.RefreshTokens
	DB.Delete(&token, tokenID)
	return GenerateRefreshToken(DB, userID)
}