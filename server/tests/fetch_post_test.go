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

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"

	token, err := auth.GenerateToken(author_id, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	ogPost := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+ogPost.Post_ID, nil)
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, second.Code)
	makeResponseAsserions(t, http.StatusOK, result_post, PostBody{Content: ogPost.Content, Tags: ogPost.Tags, Public: ogPost.Public}, author_id, second.Code)
}

func TestGetPostWithInvalidID(t *testing.T) {

	log.Println("TestGetPostWithInvalidID")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "54"
	token, err := auth.GenerateToken(author_id, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	ogPost := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+ogPost.Post_ID+"invalid", nil)
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusNotFound, second.Code)
}
