package posts_test

import (
	"encoding/json"

	"github.com/andro-kes/Blog/controllers/posts"
	"github.com/andro-kes/Blog/utils"
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/tests/helpers"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

func TestGetUserPostsHandler(t *testing.T) {
	router := SetUpTestRouter()
	router.GET("posts/all/:id", posts_controllers.GetUserPostsHandler)
	cookies := test_helpers.Login(router)

	var token string
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			token = cookie.Value
		}
	}
	assert.NotEmpty(t, token)

	claims, err := utils.ParseRefreshToken(token)
	assert.NoError(t, err)

	userID := claims.UserID
	assert.True(t, userID != 0)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"GET", 
		fmt.Sprintf("/posts/all/%d", userID),
		nil,
	)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var posts []models.Posts
	err = json.Unmarshal(w.Body.Bytes(), &posts)

	assert.Equal(t, 1, len(posts))
}