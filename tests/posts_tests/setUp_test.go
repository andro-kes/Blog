package posts_test

import (
	"github.com/andro-kes/Blog/middlewares"
	"github.com/andro-kes/Blog/models"
	

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

var DB *gorm.DB

func SetUpTestRouter() (*gin.Engine) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middlewares.DBMiddleWare())

	DB = middlewares.TestDB

	DB.Migrator().DropTable(&models.Posts{})
	DB.Migrator().DropTable(&models.Users{})
	DB.Migrator().DropTable(&models.RefreshTokens{})

	DB.Migrator().CreateTable(&models.Users{})
	DB.Migrator().CreateTable(&models.Posts{})
	DB.Migrator().CreateTable(&models.RefreshTokens{})

	password, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), 16)

	user := models.Users{
		Email: "testemail",
		Password: string(password),
		Role: "user",
		IsOauth: false,
	}
	DB.Create(&user)

	post := models.Posts{
		UserID: uint(1),
		Text: "testtext",
	}
	DB.Create(&post)

	return router
}