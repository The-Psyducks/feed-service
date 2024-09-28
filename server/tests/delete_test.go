package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"server/src/auth"
	"server/src/models"
	"server/src/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletePost(t *testing.T) {

	log.Println("TestDeletePost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
	postBody := PostBody{Content: "content", Tags: []string{"tag1", "tag2"}, Public: true}
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

	assert.Equal(t, err, nil)
	makeResponseAsserions(t, http.StatusCreated, result, postBody, author_id, first.Code)

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+result.Post_ID, nil)
	addAuthorization(deletePost, token)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNoContent, third.Code)

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+result.Post_ID, nil)
	addAuthorization(getPost, token)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPost)
	assert.Equal(t, http.StatusNotFound, fourth.Code)
}

func TestDeleteUnexistentPost(t *testing.T) {

	log.Println("TestDeleteUnexistentPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
	postBody := PostBody{Content: "content", Tags: []string{"tag1", "tag2"}, Public: true}
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

	assert.Equal(t, err, nil)
	makeResponseAsserions(t, http.StatusCreated, result, postBody, author_id, first.Code)

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+result.Post_ID+"invalid", nil)
	addAuthorization(deletePost, token)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNotFound, third.Code)
}
