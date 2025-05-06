package posts_controllers

import (
	"github.com/andro-kes/Blog/controllers/helpers"
	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"

	"log"
)

func UpdatePostHandler(c *gin.Context) {
	DB := controllers_helpers.Connect_db(c)
	if DB == nil {
		return
	}

	id := c.Param("id")
	var post models.Posts
	obj := DB.Where("id = ?", id).First(&post)
	if obj.Error != nil {
		log.Println("UpdatePost: Пост не найден")
		c.JSON(400, gin.H{"error": "Такого поста не существует"})
		return
	}

	type Content struct {
		Text string `json:"text"`
	} 
	var content Content

	if err := c.ShouldBindJSON(&content); err != nil {
		log.Println("UpdatePost: Не удалось связать со структурой")
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	obj = DB.Model(&post).Update("text", content.Text)
	if obj.Error != nil {
		log.Println("UpdatePost: Ошибка обновления поста")
		c.JSON(400, gin.H{"error": "Не удалось обновить пост"})
		return
	}

	c.JSON(200, post)
}