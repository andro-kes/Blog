package posts_controllers

import (
	"log"

	"github.com/andro-kes/Blog/controllers/helpers"
	"github.com/andro-kes/Blog/models"
	"github.com/gin-gonic/gin"
)

func GetAllPostsHandler(c *gin.Context) {
	DB := controllers_helpers.Connect_db(c)
	if DB == nil {
		return
	}

	var posts []models.Posts
	DB.Find(&posts)
	if len(posts) == 0 {
		log.Println("Постов нет")
		c.JSON(200, gin.H{"message": "постов нет"})
	} else {
		log.Printf("Было найдено %d постов", len(posts))
		c.JSON(200, posts)
	}
}