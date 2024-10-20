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

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"

	token, err := auth.GenerateToken(author_id, "username", false)

	assert.Equal(t, err, nil)

	ogPost := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	newContent := "new content"
	newTags := []string{}

	editInfo := struct {
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}{
		Content: newContent,
		Tags:    newTags,
	}

	newPostBody := PostBody{Content: newContent, Tags: ogPost.Tags, Public: ogPost.Public}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil)
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPostTags(t *testing.T) {

	log.Println("TestEditPostTags")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"

	token, err := auth.GenerateToken(author_id, "username", false)

	assert.Equal(t, err, nil)

	ogPost := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	newContent := ""
	newTags := []string{"New", "Tags"}

	editInfo := struct {
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}{
		Content: newContent,
		Tags:    newTags,
	}

	newPostBody := PostBody{Content: ogPost.Content, Tags: newTags, Public: ogPost.Public}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, second.Code)
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPostMediaURL(t *testing.T) {

	log.Println("TestEditPostMediaURL")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"

	token, err := auth.GenerateToken(author_id, "username", false)

	base_media_url := "media_url_original"

	edit_media_url := "media_url_edited"

	assert.Equal(t, err, nil)

	ogPost := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, base_media_url, r, t)

	newContent := ""
	newTags := []string{}

	editInfo := struct {
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
		MediaURL string `json:"media_url"`
	}{
		Content: newContent,
		Tags:    newTags,
		MediaURL: edit_media_url,
	}

	newPostBody := PostBody{Content: ogPost.Content, Tags: ogPost.Tags, Public: ogPost.Public, MediaURL: edit_media_url}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil)
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPost(t *testing.T) {
	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
	token, err := auth.GenerateToken(author_id, "username", false)

	assert.Equal(t, err, nil)
	
	ogPost := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	newContent := "new content"
	newTags := []string{"New", "Tags"}
	edit_media_url := "media_url_edited"

	editInfo := struct {
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
		MediaURL string `json:"media_url"`
	}{
		Content: newContent,
		Tags:    newTags,
		MediaURL: edit_media_url,
	}

	newPostBody := PostBody{Content: newContent, Tags: newTags, Public: ogPost.Public, MediaURL: edit_media_url}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, second.Code)
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}
