package test

import (
	"bytes"
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

func TestEditPostContent(t *testing.T) {

	log.Println("TestEditPostContent")

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

	newContent := "new content"
	newTags := []string{}

	editInfo := struct {
		Content  string  `json:"content"`
		Tags []string `json:"tags"`
	}{
        Content: newContent,
        Tags: newTags,
    }

	newPostBody := PostBody{Content: newContent, Tags: postBody.Tags, Public: postBody.Public}
	
	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+result.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	AddAuthorization(getPost, token)


	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil)
	MakeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPostTags(t *testing.T) {

	log.Println("TestEditPostTags")

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

	newContent := ""
	newTags := []string{"New", "Tags"}

	editInfo := struct {
		Content  string  `json:"content"`
		Tags []string `json:"tags"`
	}{
        Content: newContent,
        Tags: newTags,
    }

	newPostBody := PostBody{Content: postBody.Content, Tags: newTags, Public: postBody.Public}
	
	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+result.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	AddAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	log.Println(result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, second.Code)
	MakeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPost(t *testing.T) {
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

	newContent := "new content"
	newTags := []string{"New", "Tags"}

	editInfo := struct {
		Content  string  `json:"content"`
		Tags []string `json:"tags"`
	}{
        Content: newContent,
        Tags: newTags,
    }

	newPostBody := PostBody{Content: newContent, Tags: newTags, Public: postBody.Public}
	
	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+result.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	AddAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	log.Println(result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, second.Code)
	MakeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}