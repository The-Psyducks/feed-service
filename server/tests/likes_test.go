package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"server/src/auth"
	"server/src/models"
	"server/src/router"
	"server/src/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLikingAPost(t *testing.T) {
	log.Println("TestLikingAPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
	liker_id := "2"

	tokenLiker, err := auth.GenerateToken(liker_id, "username", true)

	assert.Equal(t, err, nil)

	post := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	getPost, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Post_ID, nil)
	addAuthorization(getPost, tokenLiker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, getPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	getPostLiked, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPostLiked, tokenLiker)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostLiked)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil)

	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, result_post.Likes, 1)
	assert.Equal(t, result_post.User_Liked, true)
}

func TestUnlikingAPost(t *testing.T) {
	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
	liker_id := "2"

	tokenLiker, err := auth.GenerateToken(liker_id, "username", true)

	assert.Equal(t, err, nil)

	post := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	getPost, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Post_ID, nil)
	addAuthorization(getPost, tokenLiker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, getPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	getPostLiked, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPostLiked, tokenLiker)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostLiked)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil)

	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, result_post.Likes, 1)
	assert.Equal(t, result_post.User_Liked, true)

	getPostUnlike, _ := http.NewRequest("DELETE", "/twitsnap/like/"+post.Post_ID, nil)
	addAuthorization(getPostUnlike, tokenLiker)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getPostUnlike)

	assert.Equal(t, http.StatusNoContent, third.Code)

	getPostUnLiked, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPostUnLiked, tokenLiker)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPostUnLiked)

	result_post_s := models.FrontPost{}

	err = json.Unmarshal(fourth.Body.Bytes(), &result_post_s)

	assert.Equal(t, err, nil)

	assert.Equal(t, http.StatusOK, fourth.Code)
	assert.Equal(t, result_post_s.Likes, 0)
	assert.Equal(t, result_post_s.User_Liked, false)
}

func TestSeeLikedTweetInFeedFollowing(t *testing.T) {
	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	liker_id := service.TEST_USER_TWO

	tokenLiker, err := auth.GenerateToken(liker_id, "username", true)
	assert.Equal(t, err, nil)

	post := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)
	time.Sleep(1 * time.Second)
	makeAndAssertPost(service.TEST_USER_THREE, "content", []string{"tag1", "tag2"}, true, "", r, t)
	time.Sleep(1 * time.Second)
	makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	getPost, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Post_ID, nil)
	addAuthorization(getPost, tokenLiker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, getPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_F+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, tokenLiker)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)

	assert.Equal(t, result.Data[2].Likes, 1)
	assert.Equal(t, result.Data[2].User_Liked, true)
}

func TestUserCanNotLikeTwice(t *testing.T) {
	log.Println("TestUserCanNotLikeTwice")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	liker_id := service.TEST_USER_TWO

	tokenLiker, err := auth.GenerateToken(liker_id, "username", true)

	assert.Equal(t, err, nil)

	post := makeAndAssertPost(author_id, "content", []string{"tag1", "tag2"}, true, "", r, t)

	getPost, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Post_ID, nil)
	addAuthorization(getPost, tokenLiker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, getPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	getPostLiked, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPostLiked, tokenLiker)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostLiked)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil)

	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, result_post.Likes, 1)
	assert.Equal(t, result_post.User_Liked, true)

	getPost2, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Post_ID, nil)
	addAuthorization(getPost2, tokenLiker)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getPost2)

	assert.Equal(t, http.StatusBadRequest, third.Code)
}
