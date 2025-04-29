package models

import (
	"gorm.io/gorm"
)

type Posts struct {
	gorm.Model
	UserID uint
	User Users `gorm:"foreignKey:UserID"`
	Text string `json:"text"`
}