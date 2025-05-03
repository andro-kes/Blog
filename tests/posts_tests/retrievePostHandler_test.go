package posts_test

import (
	"encoding/json"

	"github.com/andro-kes/Blog/controllers/posts"
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/tests/helpers"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

func TestRetrievePostHandler(t *testing.T) {
	router := SetUpTestRouter()
	router.GET("posts/:id", posts_controllers.RetrievePostHandler)
	cookies := test_helpers.Login(router)

	var existingPost models.Posts
	obj := DB.First(&existingPost)
	assert.NoError(t, obj.Error)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/posts/%d", existingPost.ID)
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response gin.H
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, existingPost.Text, response["text"])
}