package models

import (
	"gorm.io/gorm"
)

type Posts struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User Users `gorm:"foreignKey:UserID"`
	Text string `json:"text"`
}