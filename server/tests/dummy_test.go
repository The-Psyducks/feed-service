package test

import (
	"log"
	"testing"
	"time"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"server/src/database"
	"server/src/router"
)

func newPostRequest(author_id string, content string, tags []string, public bool, r *gin.Engine) *httptest.ResponseRecorder {
	post := struct {
		Author_ID string `json:"author_id"`
		Content  string  `json:"content"`
		Tags []string `json:"tags"`
		Public bool `json:"public"`
	}{
        Author_ID: author_id,
        Content: content,
        Tags: tags,
		Public: public,
    }
	
	marshalledData, _ := json.Marshal(post)
	req, _ := http.NewRequest("POST", "/twitsnap", bytes.NewReader(marshalledData))

	req.Header.Add("content-type", "application/json")


	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	return first
}
 

func TestNewPost(t *testing.T) {

    db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"tag1", "tag2"}
	public := true

    first := newPostRequest(author_id, content,tags,public, r)

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
	assert.Equal(t, http.StatusCreated, first.Code)
	assert.Equal(t, result.Post.Content, content)
	assert.Equal(t, result.Post.Author_ID, author_id)
	assert.Equal(t, result.Post.Tags, tags)
	assert.Equal(t, result.Post.Public, public)
}

func TestGetPost(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"tag1", "tag2"}
	public := true

    first := newPostRequest(author_id, content,tags,public, r)

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
}

func TestDeletePost(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"tag1", "tag2"}
	public := true

    first := newPostRequest(author_id, content,tags,public, r)

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

func TestGetFeed(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"pencil", "kiwi"}
	public := true

    _ = newPostRequest(author_id, content,tags,public, r)
	author_id_second := "2"
	content_second  := "second twitsnap content"
	tags_second  := []string{"apple", "pie"}
	public_second  := false

	time.Sleep(1 * time.Second)

    _  = newPostRequest(author_id_second, content_second,tags_second,public_second, r)

    result := struct {
		Posts []struct {
			Posts_ID   string    `bson:"post_id"`
			Content   string    `bson:"content"`
			Author_ID string    `bson:"author_id"`
			Time      time.Time `bson:"time"`
			Public   bool    `bson:"public"`
			Tags     []string  `bson:"tags"`
		}
	}{}
	
	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed", nil)
	
	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)
	
	err := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	log.Println(result)
	
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)
	assert.Equal(t, result.Posts[0].Content, content_second)
	assert.Equal(t, result.Posts[0].Author_ID, author_id_second)
	assert.Equal(t, result.Posts[0].Tags, tags_second)
	assert.Equal(t, result.Posts[0].Public, public_second)
	assert.Equal(t, result.Posts[1].Content, content)
	assert.Equal(t, result.Posts[1].Author_ID, author_id)
	assert.Equal(t, result.Posts[1].Tags, tags)
	assert.Equal(t, result.Posts[1].Public, public)
}