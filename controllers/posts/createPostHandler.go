package posts_controllers

import (
	"github.com/andro-kes/Blog/controllers/helpers"
	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	DB := controllers_helpers.Connect_db(c)
	if DB == nil {
		return
	}

	var post models.Posts
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": "Невалидные данные для создания поста"})
		return
	}

	DB.Create(&post)

	c.JSON(201, gin.H{
		"message": "Пост был успешно создан",
		"text": post.Text,
	})
}