package controllers_test

import (
	"bytes"

	"github.com/andro-kes/Blog/controllers/users"
	"github.com/andro-kes/Blog/middlewares"
	"github.com/andro-kes/Blog/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

	DB.Migrator().CreateTable(&models.Users{})
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

func TestLoginHandler(t *testing.T) {
	router := SetUpTestRouter()
	router.POST("users/login", controllers.LoginHandler)
	user := models.Users{
		Email: "testemail",
		Password: "testpassword",
	}
	jsonUser, err := json.Marshal(user)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonUser))
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response gin.H
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, response["email"])
}