package test

import (
	"log"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"server/src/auth"
	"server/src/models"
	"server/src/router"
)

func TestGetPostWithValidID(t *testing.T) {

	log.Println("TestGetPostWithValidID")

	db := ConnectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
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


	getPost, _ := http.NewRequest("GET", "/twitsnap/"+result.Post_ID, nil)
	AddAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	log.Println(result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, second.Code)
	MakeResponseAsserions(t, http.StatusOK, result_post, postBody, author_id, second.Code)
}

func TestGetPostWithInvalidID(t *testing.T) {

	log.Println("TestGetPostWithInvalidID")

	db := ConnectToDatabase()

	r := router.CreateRouter(db)

	author_id := "54"
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

	MakeResponseAsserions(t, http.StatusCreated, result, postBody, author_id, first.Code)

	assert.Equal(t, err, nil)

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+result.Post_ID+"invalid", nil)
	AddAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusNotFound, second.Code)
}
