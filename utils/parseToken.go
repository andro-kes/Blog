package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/config"
)

func ParseToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if ok == false {
		return nil, err
	}

	return claims, err
}