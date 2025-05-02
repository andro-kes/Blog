package users_test

import (
	"github.com/andro-kes/Blog/tests/helpers"
	"github.com/andro-kes/Blog/controllers/users"
	"github.com/andro-kes/Blog/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutHandler(t *testing.T) {
	router := SetUpTestRouter()
	cookies := test_helpers.Login(router)

	router.POST("users/logout", users_controllers.LogoutHandler)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/users/logout", nil)
	assert.NoError(t, err)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	c := gin.CreateTestContextOnly(w, router)
	c.Request = req
	
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	cookies = w.Result().Cookies()
	var (
		token string
		refresh_token string
	)
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			token = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			refresh_token = cookie.Value
		}
	}
	assert.Empty(t, token)
	assert.Empty(t, refresh_token)

	var RefreshToken models.RefreshTokens
	DB.Where("user_id = ?", uint(1)).First(&RefreshToken)
	assert.Equal(t, "", RefreshToken.Token)
}