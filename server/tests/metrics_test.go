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

func TestMetricsLikes(t *testing.T) {
	log.Println("TestMetricsLikes")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	time_init := time.Now().Format(time.RFC3339)

	author_id := service.TEST_USER_ONE
	liker_ids := []string{service.TEST_USER_TWO, service.TEST_USER_THREE}

	tokenAuthor, err := auth.GenerateToken(author_id, service.TEST_USER_ONE_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tokenLiker, err := auth.GenerateToken(liker_ids[0], service.TEST_USER_TWO_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tokenLiker1, err := auth.GenerateToken(liker_ids[1], service.TEST_USER_THREE_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	time.Sleep(1 * time.Second)

	post := makeAndAssertPost(author_id, "content "+"#"+tags[0]+" #"+tags[1], tags, []string{}, true, "", r, t)

	getPost, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Original_Post_ID, nil)
	addAuthorization(getPost, tokenLiker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, getPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	getPost2, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Original_Post_ID, nil)
	addAuthorization(getPost2, tokenLiker1)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost2)

	assert.Equal(t, http.StatusNoContent, second.Code)

	time.Sleep(1 * time.Second)

	endTime := time.Now().Format(time.RFC3339)

	getPostLiked, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getPostLiked, tokenLiker)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPostLiked)

	result := models.FrontPost{}

	err = json.Unmarshal(fourth.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")

	// log.Println(result)

	assert.Equal(t, http.StatusOK, fourth.Code, "Status should be 200")
	assert.Equal(t, result.Likes, 2, "Post should have 2 likes")
	assert.Equal(t, result.User_Liked, true, "User should have liked the post")

	getMetrics, _ := http.NewRequest("GET", "/twitsnap/metrics?time="+time_init+"&end_time="+endTime, nil)
	addAuthorization(getMetrics, tokenAuthor)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getMetrics)

	result_post := models.UserMetrics{}

	err = json.Unmarshal(third.Body.Bytes(), &result_post)

	log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, third.Code, "Status should be 200")
	assert.Equal(t, result_post.Likes, 2, "Post should have 2 likes")
	assert.Equal(t, result_post.Retweets, 0, "Post should have 0 retweets")
	assert.Equal(t, result_post.Posts, 1, "There should be 1 post")
}

func TestMetricsRetweets(t *testing.T) {
	log.Println("TestMetricsRetweets")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	time_init := time.Now().Format(time.RFC3339)

	author_id := service.TEST_USER_ONE
	rt_ids := []string{service.TEST_USER_TWO, service.TEST_USER_THREE}

	tokenAuthor, err := auth.GenerateToken(author_id, "username", true)
	assert.Equal(t, err, nil, "Error should be nil")

	tokenRT, err := auth.GenerateToken(rt_ids[0], service.TEST_USER_TWO_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tokenRT1, err := auth.GenerateToken(rt_ids[1], service.TEST_USER_THREE_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	time.Sleep(1 * time.Second)

	post := makeAndAssertPost(author_id, "content "+"#"+tags[0]+" #"+tags[1], tags, []string{}, true, "", r, t)

	retweetAPost(post, service.TEST_USER_TWO_USERNAME, tokenRT, r, t)

	getPostRetweet, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getPostRetweet, tokenAuthor)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPostRetweet)

	result := models.FrontPost{}

	err = json.Unmarshal(fourth.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")

	retweetAPost(result, service.TEST_USER_THREE_USERNAME, tokenRT1, r, t)

	time.Sleep(1 * time.Second)

	endTime := time.Now().Format(time.RFC3339)

	getMetrics, _ := http.NewRequest("GET", "/twitsnap/metrics?time="+time_init+"&end_time="+endTime, nil)
	addAuthorization(getMetrics, tokenAuthor)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getMetrics)

	result_post := models.UserMetrics{}

	err = json.Unmarshal(third.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, third.Code, "Status should be 200")
	assert.Equal(t, result_post.Likes, 0, "Post should have 2 likes")
	assert.Equal(t, result_post.Retweets, 2, "Post should have 0 retweets")
	assert.Equal(t, result_post.Posts, 1, "There should be 1 post")
}

func TestLikesAndRetweets(t *testing.T) {
	log.Println("TestLikesAndRetweets")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	time_init := time.Now().Format(time.RFC3339)

	author_id := service.TEST_USER_ONE
	liker_ids := []string{service.TEST_USER_TWO, service.TEST_USER_THREE}

	tokenAuthor, err := auth.GenerateToken(author_id, service.TEST_USER_ONE_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tokenLiker, err := auth.GenerateToken(liker_ids[0], service.TEST_USER_TWO_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tokenLiker1, err := auth.GenerateToken(liker_ids[1], service.TEST_USER_THREE_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	time.Sleep(1 * time.Second)

	post := makeAndAssertPost(author_id, "content "+"#"+tags[0]+" #"+tags[1], tags, []string{}, true, "", r, t)

	getPost, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Original_Post_ID, nil)
	addAuthorization(getPost, tokenLiker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, getPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	getPost2, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Original_Post_ID, nil)
	addAuthorization(getPost2, tokenLiker1)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPost2)

	assert.Equal(t, http.StatusNoContent, second.Code)

	getPostLiked, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getPostLiked, tokenAuthor)

	fifth := httptest.NewRecorder()
	r.ServeHTTP(fifth, getPostLiked)

	res := models.FrontPost{}

	err = json.Unmarshal(fifth.Body.Bytes(), &res)

	assert.Equal(t, err, nil, "Error should be nil")

	retweetAPost(res, service.TEST_USER_TWO_USERNAME, tokenLiker, r, t)

	getPostRetweet, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getPostRetweet, tokenAuthor)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPostRetweet)

	result := models.FrontPost{}

	err = json.Unmarshal(fourth.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")

	retweetAPost(result, service.TEST_USER_THREE_USERNAME, tokenLiker1, r, t)

	time.Sleep(1 * time.Second)

	endTime := time.Now().Format(time.RFC3339)

	getMetrics, _ := http.NewRequest("GET", "/twitsnap/metrics?time="+time_init+"&end_time="+endTime, nil)
	addAuthorization(getMetrics, tokenAuthor)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getMetrics)

	result_post := models.UserMetrics{}

	err = json.Unmarshal(third.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, third.Code, "Status should be 200")
	assert.Equal(t, result_post.Likes, 2, "Post should have 2 likes")
	assert.Equal(t, result_post.Retweets, 2, "Post should have 2 retweets")
	assert.Equal(t, result_post.Posts, 1, "There should be 1 post")
}
