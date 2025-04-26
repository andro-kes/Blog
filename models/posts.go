package models

import (
	"gorm.io/gorm"
)

type Posts struct {
	gorm.Model
	UserID uint
	User Users `gorm:"foreignKey:UserID"`
	Description string `json:"description"`
	PictureData []byte `gorm:"type:longblob"`
}