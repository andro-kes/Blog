package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CompareHashPasswords(userPassword, existingPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(existingPassword))
}