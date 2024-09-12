package test

import (
	"log"
	"testing"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	// "github.com/mjarkk/mongomock"
	"github.com/stretchr/testify/assert"

	"server/src/database"
	"server/src/router"
)

func TestDeletePost(t *testing.T) {
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

	log.Println(result)

	assert.Equal(t, err, nil)

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+result.Post.Post_ID, nil)

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
	assert.Equal(t, result_post.Post.Tags, tags)
	assert.Equal(t, result_post.Post.Public, public)

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+result.Post.Post_ID, nil)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNoContent, third.Code)

	getPost, _ = http.NewRequest("GET", "/twitsnap/"+result.Post.Post_ID, nil)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPost)

	assert.Equal(t, http.StatusNotFound, fourth.Code)
}

func TestDeleteUnexistentPost(t *testing.T) {
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

	log.Println(result)

	assert.Equal(t, err, nil)

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+result.Post.Post_ID, nil)

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
	assert.Equal(t, result_post.Post.Tags, tags)
	assert.Equal(t, result_post.Post.Public, public)

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+result.Post.Post_ID+"invalid", nil)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNotFound, third.Code)
}