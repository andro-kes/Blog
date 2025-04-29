package controllers

import (
	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Пользователь вышел из системы"})
}