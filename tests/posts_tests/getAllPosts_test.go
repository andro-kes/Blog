package posts_test

import (
	"encoding/json"

	"github.com/andro-kes/Blog/controllers/posts"
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/tests/helpers"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllPostsHandler(t *testing.T) {
	router := SetUpTestRouter()
	router.GET("posts/all", posts_controllers.GetAllPostsHandler)
	cookies := test_helpers.Login(router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/posts/all", nil)
	assert.NoError(t, err)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response []models.Posts
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 1, len(response))
}