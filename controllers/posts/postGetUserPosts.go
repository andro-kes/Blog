package posts_controllers

import (
	"log"

	"github.com/andro-kes/Blog/controllers/helpers"
	"github.com/andro-kes/Blog/utils"
	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"
)

func GetUserPostsHandler(c *gin.Context) {
	DB := controllers_helpers.Connect_db(c)
	if DB == nil {
		return
	}

	refresh_token , err := c.Cookie("refresh_token")
	if err != nil {
		log.Println("Refresh token не найденв куках")
		c.JSON(400, "Refresh token не найден")
		return
	}

	claims, err := utils.ParseRefreshToken(refresh_token)
	if err != nil {
		log.Println("Не удалось извлечь claims")
		c.JSON(400, "Доступ запрещен")
		return
	}

	var posts []models.Posts
	DB.Where("user_id = ?", claims.UserID).Find(&posts)
	if len(posts) == 0 {
		log.Println("Нет постов")
		c.JSON(200, gin.H{"message": "У этого пользователя нет постов"})
		return
	} else {
		log.Printf("Найдено %d постов", len(posts))
		c.JSON(200, posts)
		return
	}
}