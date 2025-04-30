package controllers_test

import (
	"bytes"
	// "log"

	"github.com/andro-kes/Blog/controllers/users"
	"github.com/andro-kes/Blog/models"
	// "github.com/andro-kes/Blog/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefreshTokenHandler(t *testing.T) {
	router := SetUpTestRouter()
	router.POST("users/login", controllers.LoginHandler)
	w := httptest.NewRecorder()
	user := models.Users{
		Email: "testemail",
		Password: "testpassword",
	}
	jsonUser, err := json.Marshal(user)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonUser))
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	cookies := w.Result().Cookies()
	var (
		token string
		refresh_token string
	)
	token, refresh_token = getTokens(cookies)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refresh_token)
	
	var RefreshToken models.RefreshTokens
	assert.NoError(t, err)
	obj := DB.Where("token = ?", refresh_token).First(&RefreshToken)
	assert.NoError(t, obj.Error)
	assert.True(t, RefreshToken.Token != "")

	router.POST("users/refresh_token", controllers.RefreshTokenHandler)
	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/users/refresh_token", nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
		// log.Printf("%s: %s\n", cookie.Name, cookie.Value)
		req.AddCookie(cookie)
	}
	c := gin.CreateTestContextOnly(w, router)
	c.Request = req
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	cookies = w.Result().Cookies()
	token, refresh_token = getTokens(cookies)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refresh_token)
	
	var NewRefreshToken models.RefreshTokens
	obj = DB.Where("token = ?", refresh_token).First(&NewRefreshToken)
	assert.NoError(t, obj.Error)
	assert.True(t, RefreshToken.Token != "")
}

func getTokens(cookies []*http.Cookie) (string, string){
	var token, refresh_token string
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			token = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			refresh_token = cookie.Value
		}
	}
	return token, refresh_token
}