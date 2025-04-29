package utils

import (
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/config"

	"github.com/dgrijalva/jwt-go"

	"time"
)

func GenerateAccessToken(existingUser models.Users) (string, error) {
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
	return tokenString, err
}

