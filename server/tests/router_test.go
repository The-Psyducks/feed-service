package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	postErrors "server/src/all_errors"
	"server/src/auth"
	"server/src/router"
	"server/src/service"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadRouterRoute(t *testing.T) {

	log.Println("TestBadRouterRoute")

	db := connectToDatabase()

	r := router.CreateRouter(db)
	token, err := auth.GenerateToken(service.TEST_USER_INFO_ERROR, "username", true)

	getPost, _ := http.NewRequest("POST", "/twitsnap/likers/"+"bad", nil)
	addAuthorization(getPost, token)


	assert.Equal(t, err, nil, "Error should be nil")

	first := httptest.NewRecorder()
	r.ServeHTTP(first, getPost)

	result := postErrors.NotFound()

	err = json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")

	errRouterError := postErrors.NotFound()

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusNotFound, first.Code)
	assert.Equal(t, errRouterError.Detail, result.Detail)
}