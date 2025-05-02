package test_helpers

import (
	"bytes"

	"github.com/andro-kes/Blog/controllers/users"
	"github.com/andro-kes/Blog/models"

	"github.com/gin-gonic/gin"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"log"
)

func Login(router *gin.Engine) []*http.Cookie {
	router.POST("users/login", users_controllers.LoginHandler)
	w := httptest.NewRecorder()
	user := models.Users{
		Email: "testemail",
		Password: "testpassword",
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Fatalln("Ошибка login")
		return nil
	}
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonUser))
	router.ServeHTTP(w, req)

	cookies := w.Result().Cookies()
	return cookies
}