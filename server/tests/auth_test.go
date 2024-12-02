package test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"server/src/auth"
	"server/src/service"

	"net/http/httptest"

	postErrors "server/src/all_errors"
	"server/src/router"

	"net/http"
)

func TestAuthenticationToken(t *testing.T) {

	log.Println("TestAuthenticationToken")

	author_id := "1234"

	token, err := auth.GenerateToken(author_id, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.NotEqual(t, token, "", "Token should not be empty")
}

func TestAuthenticationTokenErrorHeaderRequired(t *testing.T) {

	log.Println("TestAuthenticationTokenErrorHeaderRequired")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	postBody := PostBody{Content: "", Tags: []string{"tag1", "tag2"}, Public: true, Mentions: []string{}}
	req := newPostRequest(postBody)

	_, err := auth.GenerateToken(service.TEST_USER_INFO_ERROR, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := postErrors.AuthenticationErrorHeaderRequired()
	errAuthError := postErrors.AuthenticationErrorHeaderRequired()

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusUnauthorized, first.Code)
	assert.Equal(t, errAuthError.Detail, result.Detail)
}