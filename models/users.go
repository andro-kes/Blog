package models

import (
	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"
)

type Users struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	IsOauth bool `gorm:"default:false"`
}

type Claims struct {
	Role string
	jwt.StandardClaims
}