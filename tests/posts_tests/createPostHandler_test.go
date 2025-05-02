package posts_test

import (
	"bytes"
	"encoding/json"

	"github.com/andro-kes/Blog/controllers/posts"
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/tests/helpers"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePostHandler(t *testing.T) {
	router := SetUpTestRouter()
	cookies := test_helpers.Login(router)

	post := models.Posts{
		UserID: uint(1),
		Text: "test",
	}
	jsonPost, err := json.Marshal(post)
	assert.NoError(t, err)

	router.POST("posts/create", posts_controllers.CreatePostHandler)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/posts/create", bytes.NewBuffer(jsonPost))
	assert.NoError(t, err)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var response gin.H
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, post.Text, response["text"])
}