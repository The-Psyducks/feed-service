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

func TestBlockedPost(t *testing.T) {
	log.Println("TestBlockedPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	retweeter_id := service.TEST_USER_TWO

	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, service.TEST_USER_TWO_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	blockPost, _ := http.NewRequest("POST", "/twitsnap/block/"+post.Post_ID, nil)
	addAuthorization(blockPost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, blockPost)


	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusNoContent, first.Code)


	getPostAfterRetweet, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPostAfterRetweet, tokenRetweeterer)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostAfterRetweet)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	// log.Println(result_post)
	assert.Equal(t, http.StatusNotFound, second.Code)
}

func TestBlockeedInFeedFollowing(t *testing.T) {
	log.Println("TestBlockeedInFeedFollowing")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{"tag1", "tag2"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, service.TEST_USER_ONE_USERNAME, false)

	assert.Equal(t, err, nil, "Error should be nil")

	retweeter := service.TEST_USER_TWO
	tokenRetweeterer, err := auth.GenerateToken(retweeter, service.TEST_USER_TWO_USERNAME, false)

	assert.Equal(t, err, nil, "Error should be nil")

	retweetPost, _ := http.NewRequest("POST", "/twitsnap/block/"+post1.Post_ID, nil)
	addAuthorization(retweetPost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, retweetPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	expectedPosts := []models.FrontPost{post3, post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_F+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)

}

func TestUnBlockAPost(t *testing.T) {
	log.Println("TestUnRetweetAPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1"
	retweeter_id := service.TEST_USER_TWO
	username := service.TEST_USER_TWO_USERNAME
	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, username, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	blockPost, _ := http.NewRequest("POST", "/twitsnap/block/"+post.Original_Post_ID, nil)
	addAuthorization(blockPost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, blockPost)

	getPostAfterRetweet, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getPostAfterRetweet, tokenRetweeterer)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostAfterRetweet)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusNoContent, first.Code)

	unretweetPost, _ := http.NewRequest("DELETE", "/twitsnap/block/"+post.Original_Post_ID, nil)
	addAuthorization(unretweetPost, tokenRetweeterer)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, unretweetPost)

	assert.Equal(t, http.StatusNoContent, third.Code)

	getPostAfterUnRetweet, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getPostAfterUnRetweet, tokenRetweeterer)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPostAfterUnRetweet)

	result_post_no_block := models.FrontPost{}

	err = json.Unmarshal(fourth.Body.Bytes(), &result_post_no_block)

	assert.Equal(t, err, nil, "Error should be nil")

	// log.Println(result_post)
	assert.Equal(t, http.StatusOK, fourth.Code)
	assert.Equal(t, false, result_post_no_block.User_Retweet)
	assert.Equal(t, result_post_no_block.Retweets, 0)
}

func TestBlockInFeedForyou(t *testing.T) {
	log.Println("TestBlockInFeedForyou")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{service.TEST_TAG_ONE, "tag5"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags = []string{service.TEST_TAG_TWO, "tag4"}

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content3 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags = []string{service.TEST_TAG_THREE, "tag6"}

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content2 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, service.TEST_USER_ONE_USERNAME, false)

	assert.Equal(t, err, nil, "Error should be nil")

	retweeter := service.TEST_USER_TWO
	username := service.TEST_USER_TWO_USERNAME
	tokenRetweeterer, err := auth.GenerateToken(retweeter, username, false)

	assert.Equal(t, err, nil, "Error should be nil")

	blockPost, _ := http.NewRequest("POST", "/twitsnap/block/"+post2.Original_Post_ID, nil)
	addAuthorization(blockPost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, blockPost)

	// log.Println(retweet_post)

	expectedPosts := []models.FrontPost{post3, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_Y+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestBlockInSingleRetweeterFeed(t *testing.T) {
	log.Println("TestRetweetInSingleRetweeterFeed")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{service.TEST_TAG_ONE, "tag2"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_ONE, "content2 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_ONE, "content3 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, service.TEST_USER_ONE_USERNAME, false)

	assert.Equal(t, err, nil, "Error should be nil")

	retweeter := service.TEST_USER_ONE
	username := service.TEST_USER_ONE_USERNAME
	tokenRetweeterer, err := auth.GenerateToken(retweeter, username, false)

	assert.Equal(t, err, nil, "Error should be nil")

	blockPost, _ := http.NewRequest("POST", "/twitsnap/block/"+post1.Original_Post_ID, nil)
	addAuthorization(blockPost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, blockPost)

	expectedPosts := []models.FrontPost{post3, post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_S+"&wanted_user_id="+retweeter, nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}