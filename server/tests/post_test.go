package test

import (
	"log"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	postErrors "server/src/all_errors"
	"server/src/auth"
	"server/src/models"
	"server/src/router"

	validator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {

	log.Println("TestNewPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "content #tag1 #tag2", Tags: []string{"tag1", "tag2"}, Mentions: []string{}, Public: true}
	req := newPostRequest(postBody)

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	addAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := models.FrontPost{}

	err = json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")
	makeResponseAsserions(t, http.StatusCreated, result, postBody, author_id, first.Code)
}

func TestNewPostWithMissInf(t *testing.T) {

	log.Println("TestNewPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "", Tags: []string{"tag1", "tag2"}, Public: true, Mentions: []string{}}
	req := newPostRequest(postBody)

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	addAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := postErrors.TwitSnapError{}

	err = json.Unmarshal(first.Body.Bytes(), &result)

	newPost := models.PostExpectedFormat{Content: "", Public: true}

	validate := validator.New()

	var errFMT postErrors.TwitSnapError
	if wrongFmtErr := validate.Struct(newPost); wrongFmtErr != nil {
		errFMT = postErrors.TwitSnapImportantFieldsMissing(wrongFmtErr)
	}

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, errFMT, result)
	assert.Equal(t, http.StatusBadRequest, first.Code)
}
