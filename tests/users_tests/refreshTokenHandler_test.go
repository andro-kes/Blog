package users_test

import (
	"github.com/andro-kes/Blog/controllers/users"
	"github.com/andro-kes/Blog/tests/helpers"
	"github.com/andro-kes/Blog/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefreshTokenHandler(t *testing.T) {
	router := SetUpTestRouter()
	cookies := test_helpers.Login(router)
	var (
		token string
		refresh_token string
	)
	token, refresh_token = getTokens(cookies)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refresh_token)
	
	var RefreshToken models.RefreshTokens
	obj := DB.Where("token = ?", refresh_token).First(&RefreshToken)
	assert.NoError(t, obj.Error)
	assert.True(t, RefreshToken.Token != "")

	router.POST("users/refresh_token", users_controllers.RefreshTokenHandler)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/users/refresh_token", nil)
	assert.NoError(t, err)
	for _, cookie := range cookies {
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