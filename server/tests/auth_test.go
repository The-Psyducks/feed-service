package test

import (
	"encoding/json"
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

	err = json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")

	errAuthError := postErrors.AuthenticationErrorHeaderRequired()

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusUnauthorized, first.Code)
	assert.Equal(t, errAuthError.Detail, result.Detail)
}

func TestAuthenticationErrorInvalidHeader(t *testing.T) {

	log.Println("TestAuthenticationErrorInvalidHeader")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	postBody := PostBody{Content: "", Tags: []string{"tag1", "tag2"}, Public: true, Mentions: []string{}}
	req := newPostRequest(postBody)

	_, err := auth.GenerateToken(service.TEST_USER_INFO_ERROR, "username", true)

	req.Header.Add("Authorization", "Bearer")

	assert.Equal(t, err, nil, "Error should be nil")

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := postErrors.AuthenticationErrorInvalidHeader()

	err = json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")

	errAuthError := postErrors.AuthenticationErrorInvalidHeader()

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusUnauthorized, first.Code)
	assert.Equal(t, errAuthError.Detail, result.Detail)
}

func TestAuthenticationErrorInvalidToken(t *testing.T) {

	log.Println("TestAuthenticationErrorInvalidToken")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "content #tag1 #tag2", Tags: []string{"tag1", "tag2"}, Mentions: []string{}, Public: true}
	req := newPostRequest(postBody)

	_, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	token := "Hound&Ram"

	addAuthorization(req, "rabbits")

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := postErrors.AuthenticationErrorInvalidToken("Invalid token")

	err = json.Unmarshal(first.Body.Bytes(), &result)
	
	log.Println("result: ", result)

	assert.Equal(t, err, nil, "Error should be nil")

	var errAuthError postErrors.TwitSnapError

	_, errFMT := auth.ValidateToken(token)

	if errFMT != nil {
		errAuthError = postErrors.AuthenticationErrorInvalidToken(errFMT.Error())
	}

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusUnauthorized, first.Code)
	assert.Equal(t, errAuthError.Detail, result.Detail)
}