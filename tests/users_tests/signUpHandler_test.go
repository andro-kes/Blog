package users_test

import (
	"bytes"

	"github.com/andro-kes/Blog/controllers/users"
	"github.com/andro-kes/Blog/models"
	"github.com/andro-kes/Blog/utils"

	"github.com/stretchr/testify/assert"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUpHandler(t *testing.T) {
	router := SetUpTestRouter()
	router.POST("users/signup", users_controllers.SignupHandler)
	user := models.Users{
		Email: "test",
		Password: "test",
		Role: "user",
		IsOauth: false,
	}
	jsonUser, err := json.Marshal(user)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonUser))
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var createdUser models.Users
	obj := DB.Where("email = ?", user.Email).First(&createdUser)
	assert.NoError(t, obj.Error)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Role, createdUser.Role)
	assert.Equal(t, user.IsOauth, createdUser.IsOauth)

	err = utils.CompareHashPasswords(user.Password, createdUser.Password)
	assert.NoError(t, err)
}