package controllers_test

import (
	

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"net/http/httptest"
	"testing"
)

func SetUpTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestLoginHandler(t *testing.T) {
	router := SetUpTestRouter()
	
}