package test

import (
	"log"
	"testing"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"server/src/database"
	"server/src/router"
)


func TestGetFeed(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"pencil", "kiwi"}
	public := true

    _ = NewPostRequest(author_id, content,tags,public, r)
	author_id_second := "2"
	content_second  := "second twitsnap content"
	tags_second  := []string{"apple", "pie"}
	public_second  := false

	time.Sleep(1 * time.Second)

    _  = NewPostRequest(author_id_second, content_second,tags_second,public_second, r)

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
	
	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?following="+author_id+"&following="+author_id_second, nil)
	
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

func TestGetFeedWithFollowingThatHaveNotPosted(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"pencil", "kiwi"}
	public := true

    _ = NewPostRequest(author_id, content,tags,public, r)
	author_id_second := "2"
	content_second  := "second twitsnap content"
	tags_second  := []string{"apple", "pie"}
	public_second  := false

	time.Sleep(1 * time.Second)

    _  = NewPostRequest(author_id_second, content_second,tags_second,public_second, r)

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

	following := []string{"3", "4"}
	
	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?following="+following[0]+"&following="+following[1], nil)
	
	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)
	
	err := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	log.Println(result)
	
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)
	assert.Equal(t, len(result.Posts), 0)

}

func TestGetInterests(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"pencil", "kiwi"}
	public := true

    _ = NewPostRequest(author_id, content,tags,public, r)
	author_id_second := "2"
	content_second  := "second twitsnap content"
	tags_second  := []string{"apple", "pie"}
	public_second  := true

	time.Sleep(1 * time.Second)

    _  = NewPostRequest(author_id_second, content_second,tags_second,public_second, r)

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
	
	getFeed, _ := http.NewRequest("GET", "/twitsnap/hashtags?tags=apple", nil)
	
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
	assert.Equal(t, len(result.Posts), 1)
}

func TestGetWordSearch(t *testing.T) {
	db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"pencil", "kiwi"}
	public := true

    _ = NewPostRequest(author_id, content,tags,public, r)
	author_id_second := "2"
	content_second  := "the best apple tart recipe with amazing caramel"
	tags_second  := []string{"apple", "pie"}
	public_second  := true

	time.Sleep(1 * time.Second)

    _  = NewPostRequest(author_id_second, content_second,tags_second,public_second, r)

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
	
	getFeed, _ := http.NewRequest("GET", "/twitsnap/wordsearch?words=caramel story", nil)
	
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
	assert.Equal(t, len(result.Posts), 1)
}