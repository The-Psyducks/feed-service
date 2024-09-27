package test

import (
	"log"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	
	"github.com/stretchr/testify/assert"
	validator "github.com/go-playground/validator/v10"
	postErrors "server/src/all_errors"
	"server/src/auth"
	"server/src/models"
	"server/src/router"
)

func TestNewPost(t *testing.T) {

	log.Println("TestNewPost")

	db := ConnectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "content", Tags: []string{"tag1", "tag2"}, Public: true}
	req := NewPostRequest(postBody, r)

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	AddAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := models.FrontPost{}

	err = json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil)
	MakeResponseAsserions(t, http.StatusCreated, result, postBody, author_id, first.Code)
}

func TestNewPostWithMissInf(t *testing.T) {

	log.Println("TestNewPost")

	db := ConnectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "", Tags: []string{"tag1", "tag2"}, Public: true}
	req := NewPostRequest(postBody, r)

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	AddAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := postErrors.TwitSnapError{}

	err = json.Unmarshal(first.Body.Bytes(), &result)

	newPost := models.PostExpectedFormat{Content: "", Tags: []string{"tag1", "tag2"}, Public: true}
	validate := validator.New()
	var errFMT postErrors.TwitSnapError
	if wrongFmtErr := validate.Struct(newPost); wrongFmtErr != nil {
		errFMT = postErrors.TwitSnapImportantFieldsMissing(wrongFmtErr)
	}

	log.Println("errFMT: ", errFMT)
	log.Println("err: ", err)

	assert.Equal(t, err, nil)
	assert.Equal(t, errFMT, result)
	assert.Equal(t, http.StatusBadRequest, first.Code)
}