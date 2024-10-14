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



func TestRetweetAPost(t *testing.T) {
	log.Println("TestRetweetAPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
	retweeter_id := "2"

	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, "username", true)
	assert.Equal(t, err, nil)
	post := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, r, t)

	getPost, _ := http.NewRequest("POST", "/twitsnap/retweet/"+post.Post_ID, nil)
	addAuthorization(getPost, tokenRetweeterer)
	first := httptest.NewRecorder()
	r.ServeHTTP(first, getPost)
	retweet_post := models.FrontPost{}
	err = json.Unmarshal(first.Body.Bytes(), &retweet_post)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusCreated, first.Code)

	getPostAfterRetweet, _ := http.NewRequest("GET", "/twitsnap/"+ post.Post_ID, nil)
	addAuthorization(getPostAfterRetweet, tokenRetweeterer)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostAfterRetweet)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil)

	log.Println(result_post)
	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, true, result_post.UserRetweet)
	assert.Equal(t, result_post.Retweets, 1)

	getPostRetweeted, _ := http.NewRequest("GET", "/twitsnap/"+ retweet_post.Post_ID, nil)
	addAuthorization(getPostRetweeted, tokenRetweeterer)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getPostRetweeted)
		
	retweet_result := models.FrontPost{}

	err = json.Unmarshal(third.Body.Bytes(), &retweet_result)

	// log.Println(retweet_result)

	assert.Equal(t, err, nil)

	assert.Equal(t, http.StatusOK, third.Code)
	assert.Equal(t, true, retweet_result.UserRetweet)
}