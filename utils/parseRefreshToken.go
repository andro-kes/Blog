package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/config"
)

func ParseRefreshToken(tokenString string) (*models.RefreshTokensClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.RefreshTokensClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.RefreshTokensClaims)
	if ok == false {
		return nil, err
	}

	return claims, err
}