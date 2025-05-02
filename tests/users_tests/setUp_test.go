package users_test

import (
	"github.com/andro-kes/Blog/middlewares"
	"github.com/andro-kes/Blog/models"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"log"
)

var DB *gorm.DB

func SetUpTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middlewares.DBMiddleWare())
	DB = middlewares.TestDB

	DB.Migrator().DropTable(&models.Users{})
	DB.Migrator().DropTable(&models.RefreshTokens{})
	DB.Migrator().DropTable(&models.RefreshTokens{})

	DB.Migrator().CreateTable(&models.Users{})
	DB.Migrator().CreateTable(&models.RefreshTokens{})
	DB.Migrator().CreateTable(&models.RefreshTokens{})

	password, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), 16)
	user := models.Users{
		Email: "testemail",
		Password: string(password),
		Role: "testrole",
		IsOauth: false,
	}
	DB.Create(&user)
	log.Println(user.ID)
	return router
}