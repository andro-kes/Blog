package posts_controllers

import (
	"github.com/andro-kes/Blog/controllers/helpers"
	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"

	"log"
)

func RetrievePostHandler(c *gin.Context) {
	DB := controllers_helpers.Connect_db(c)
	if DB == nil {
		return
	}

	id := c.Param("id")
	var post models.Posts
	obj := DB.Where("id = ?", id).First(&post)
	if obj.Error != nil {
		log.Printf("Не удалось получить пост по айди: %s\n", id)
		c.JSON(400, gin.H{"error": "Пост не найден"})
		return
	}

	c.JSON(200, post)
}