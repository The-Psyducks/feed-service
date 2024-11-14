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
	"server/src/service"
)

func TestEditPostContent(t *testing.T) {

	log.Println("TestEditPostContent")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE

	token, err := auth.GenerateToken(author_id, "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	ogPost := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	newContent := "new content"

	editInfo := struct {
		Content string `json:"content"`
	}{
		Content: newContent,
	}

	newPostBody := PostBody{Content: newContent, Tags: []string{}, Public: ogPost.Public, Mentions: []string{}}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPostTags(t *testing.T) {

	log.Println("TestEditPostTags")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"

	token, err := auth.GenerateToken(author_id, "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	ogPost := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	newTags := []string{"New", "Tags"}

	newContent := "content " + "#" + newTags[0] + " #" + newTags[1]

	editInfo := struct {
		Content string `json:"content"`
	}{
		Content: newContent,
	}

	newPostBody := PostBody{Content: newContent, Tags: newTags, Public: ogPost.Public, Mentions: []string{}}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
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

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	ogPost := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, base_media_url, r, t)

	newMedia := models.MediaInfo{Media_URL: edit_media_url, Media_Type: "IMAGE"}

	editInfo := struct {
		Media_Info models.MediaInfo `json:"media_info"`
	}{
		Media_Info: newMedia,
	}

	newPostBody := PostBody{Content: ogPost.Content, Tags: tags, Public: ogPost.Public, MediaInfo: newMedia, Mentions: []string{}}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPostPublicToPrivate(t *testing.T) {

	log.Println("TestEditPostMediaURL")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"

	token, err := auth.GenerateToken(author_id, "username", false)

	public := true

	newPublic := false

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	ogPost := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, public, "base_media_url", r, t)

	editInfo := struct {
		Public bool `json:"public"`
	}{
		Public: newPublic,
	}

	newPostBody := PostBody{Content: ogPost.Content, Tags: tags, Public: newPublic, MediaInfo: ogPost.Media_Info, Mentions: []string{}}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPostPrivateToPublic(t *testing.T) {

	log.Println("TestEditPostMediaURL")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"

	token, err := auth.GenerateToken(author_id, "username", false)

	public := false

	newPublic := true

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	ogPost := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, public, "base_media_url", r, t)

	editInfo := struct {
		Public bool `json:"public"`
	}{
		Public: newPublic,
	}

	newPostBody := PostBody{Content: ogPost.Content, Tags: tags, Public: newPublic, MediaInfo: ogPost.Media_Info, Mentions: []string{}}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditMentions(t *testing.T) {

	log.Println("TestEditMentions")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE

	token, err := auth.GenerateToken(author_id, "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	mentions := []string{"user1", "user2"}

	ogPost := makeAndAssertPost(author_id, "content " + "@" + mentions[0] + " @" + mentions[1], []string{}, mentions, true, "", r, t)

	newMentions := []string{"user3"}

	newContent := "content " + "@" + newMentions[0]

	editInfo := struct {
		Content string `json:"content"`
	}{
		Content: newContent,
	}

	newPostBody := PostBody{Content: newContent, Tags: []string{}, Public: ogPost.Public, Mentions: newMentions}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}

func TestEditPostTagsAndMentions(t *testing.T) {
	log.Println("TestEditPostTagsAndMentions")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE

	token, err := auth.GenerateToken(author_id, "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	mentions := []string{"user1", "user2"}

	tags := []string{"tag1", "tag2"}

	ogPost := makeAndAssertPost(author_id, "content " + "#" + tags[0]  + " @" + mentions[0] + " #" + tags[1] + " @" + mentions[1], tags, mentions, true, "", r, t)

	newMentions := []string{"user3", "user4"}
	newTags := []string{"New", "Tags"}

	newContent := "content " + "#" + newTags[0]  + " @" + newMentions[0] + " #" + newTags[1] + " @" + newMentions[1]

	editInfo := struct {
		Content string `json:"content"`
	}{
		Content: newContent,
	}

	newPostBody := PostBody{Content: newContent, Tags: newTags, Public: ogPost.Public, Mentions: newMentions}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)

}

func TestEditPost(t *testing.T) {
	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
	token, err := auth.GenerateToken(author_id, "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	ogPost := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	newTags := []string{"New", "Tags"}
	newContent := "new content " + "#" + newTags[0] + " #" + newTags[1]
	edit_media_url := "media_url_edited"
	pubic := false

	newMedia := models.MediaInfo{Media_URL: edit_media_url, Media_Type: "IMAGE"}

	editInfo := struct {
		Content  string   `json:"content"`
		Tags     []string `json:"tags"`
		MediaURL models.MediaInfo   `json:"media_info"`
		Public   bool     `json:"public"`
	}{
		Content:  newContent,
		Tags:     newTags,
		MediaURL: newMedia,
		Public:   pubic,
	}

	newPostBody := PostBody{Content: newContent, Tags: newTags, Public: pubic, MediaInfo: newMedia, Mentions: []string{}}

	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+ogPost.Post_ID, bytes.NewBuffer(marshalledData))
	getPost.Header.Add("content-type", "application/json")
	addAuthorization(getPost, token)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
	makeResponseAsserions(t, http.StatusOK, result_post, newPostBody, author_id, second.Code)
}
