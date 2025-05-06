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

func TestUpdatePostHandler(t *testing.T) {
	router := SetUpTestRouter()
	router.PATCH("posts/update/:id", posts_controllers.UpdatePostHandler)
	cookies := test_helpers.Login(router)

	w := httptest.NewRecorder()

	content := gin.H{"text": "new text"}
	jsonContent, err := json.Marshal(content)
	assert.NoError(t, err)

	req, err := http.NewRequest("PATCH", "/posts/update/1", bytes.NewBuffer(jsonContent))
	assert.NoError(t, err)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response models.Posts
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, content["text"], response.Text)
}