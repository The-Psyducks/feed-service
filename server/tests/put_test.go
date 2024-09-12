package test

import (
	"log"
	"time"
	"bytes"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"server/src/database"
	"server/src/router"
)

func TestEditPostContent(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"tag1", "tag2"}
	public := true

    first := NewPostRequest(author_id, content,tags,public, r)

    result := struct {
		Post struct {
			Post_ID   string    `bson:"post_id"`
			Content   string    `bson:"content"`
			Author_ID string    `bson:"author_id"`
			Time      time.Time `bson:"time"`
			Public   bool    `bson:"public"`
			Tags     []string  `bson:"tags"`
		}
	}{}

	err := json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil)

	newContent := "new content"
	newTags := []string{}

	editInfo := struct {
		Content  string  `json:"content"`
		Tags []string `json:"tags"`
	}{
        Content: newContent,
        Tags: newTags,
    }
	
	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+result.Post.Post_ID, bytes.NewBuffer(marshalledData))

	getPost.Header.Add("content-type", "application/json")

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := struct {
		Post struct {
			Post_ID   string    `bson:"post_id"`
			Content   string    `bson:"content"`
			Author_ID string    `bson:"author_id"`
			Time      time.Time `bson:"time"`
			Public   bool    `bson:"public"`
			Tags     []string  `bson:"tags"`
		}
	}{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	log.Println(result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, result_post.Post.Content, newContent)
	assert.Equal(t, result_post.Post.Author_ID, author_id)
	assert.Equal(t, result_post.Post.Public, public)
	assert.Equal(t, result_post.Post.Tags, tags)
}

func TestEditPostTags(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"tag1", "tag2"}
	public := true

    first := NewPostRequest(author_id, content,tags,public, r)

    result := struct {
		Post struct {
			Post_ID   string    `bson:"post_id"`
			Content   string    `bson:"content"`
			Author_ID string    `bson:"author_id"`
			Time      time.Time `bson:"time"`
			Public   bool    `bson:"public"`
			Tags     []string  `bson:"tags"`
		}
	}{}

	err := json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil)

	newContent := ""
	newTags := []string{"New", "Tags"}

	editInfo := struct {
		Content  string  `json:"content"`
		Tags []string `json:"tags"`
	}{
        Content: newContent,
        Tags: newTags,
    }
	
	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+result.Post.Post_ID, bytes.NewBuffer(marshalledData))

	getPost.Header.Add("content-type", "application/json")

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := struct {
		Post struct {
			Post_ID   string    `bson:"post_id"`
			Content   string    `bson:"content"`
			Author_ID string    `bson:"author_id"`
			Time      time.Time `bson:"time"`
			Public   bool    `bson:"public"`
			Tags     []string  `bson:"tags"`
		}
	}{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	log.Println(result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, result_post.Post.Content, content)
	assert.Equal(t, result_post.Post.Author_ID, author_id)
	assert.Equal(t, result_post.Post.Public, public)
	assert.Equal(t, result_post.Post.Tags, newTags)
}

func TestEditPost(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"tag1", "tag2"}
	public := true

    first := NewPostRequest(author_id, content,tags,public, r)

    result := struct {
		Post struct {
			Post_ID   string    `bson:"post_id"`
			Content   string    `bson:"content"`
			Author_ID string    `bson:"author_id"`
			Time      time.Time `bson:"time"`
			Public   bool    `bson:"public"`
			Tags     []string  `bson:"tags"`
		}
	}{}

	err := json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil)

	newContent := "new content"
	newTags := []string{"New", "Tags"}

	editInfo := struct {
		Content  string  `json:"content"`
		Tags []string `json:"tags"`
	}{
        Content: newContent,
        Tags: newTags,
    }
	
	marshalledData, _ := json.Marshal(editInfo)

	getPost, _ := http.NewRequest("PUT", "/twitsnap/edit/"+result.Post.Post_ID, bytes.NewBuffer(marshalledData))

	getPost.Header.Add("content-type", "application/json")

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost)

	result_post := struct {
		Post struct {
			Post_ID   string    `bson:"post_id"`
			Content   string    `bson:"content"`
			Author_ID string    `bson:"author_id"`
			Time      time.Time `bson:"time"`
			Public   bool    `bson:"public"`
			Tags     []string  `bson:"tags"`
		}
	}{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	log.Println(result_post)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, result_post.Post.Content, newContent)
	assert.Equal(t, result_post.Post.Author_ID, author_id)
	assert.Equal(t, result_post.Post.Public, public)
	assert.Equal(t, result_post.Post.Tags, newTags)
}