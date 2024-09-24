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

func TestDeletePost(t *testing.T) {

	log.Println("TestDeletePost")

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


	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+result.Post_ID, nil)
	AddAuthorization(deletePost, token)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNoContent, third.Code)

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+result.Post_ID, nil)
	AddAuthorization(getPost, token)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPost)
	assert.Equal(t, http.StatusNotFound, fourth.Code)
}

func TestDeleteUnexistentPost(t *testing.T) {

	log.Println("TestDeleteUnexistentPost")

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

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+result.Post_ID+"invalid", nil)
	AddAuthorization(deletePost, token)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNotFound, third.Code)
}