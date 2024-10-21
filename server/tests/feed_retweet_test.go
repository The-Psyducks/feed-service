package test

import (
	"log"
	"strconv"
	"testing"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"server/src/auth"
	"server/src/models"
	"server/src/router"
	"server/src/service"
)

const (
	FEED_TYPE_R = "retweet"
)

func TestGetFeedRetweet(t *testing.T) {
	log.Println("TestGetFeedRetweet")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{"tag1", "tag2"}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{"tag3", "tag4"}, true, "", r, t)

	time.Sleep(1 * time.Second)

	retweeter_id := service.TEST_USER_THREE
	username := service.TEST_USER_THREE_USERNAME
	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, username, true)
	assert.Equal(t, err, nil, "Error should be nil")

	retweet_post1 := retweetAPost(post1, username, tokenRetweeterer, r, t)
	retweet_post2 := retweetAPost(post2, username, tokenRetweeterer, r, t)

	makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{"tag5", "tag6"}, true, "", r, t)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{retweet_post2, retweet_post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_R+"&wanted_user_id="+retweeter_id, nil)
	addAuthorization(getFeed, tokenRetweeterer)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)

	for i := range result.Data {
		checkRetweetPost(result.Data[i], username, t)
	}
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestGetFeedRetweetNotFollowing(t *testing.T) {
	log.Println("TestGetFeedRetweetNotFollowing")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{"tag1", "tag2"}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{"tag3", "tag4"}, true, "", r, t)

	time.Sleep(1 * time.Second)

	retweeter_id := service.TEST_NOT_FOLLOWING_ID
	username := service.TEST_NOT_FOLLOWING_USERNAME
	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, username, true)
	assert.Equal(t, err, nil, "Error should be nil")

	retweet_post1 := retweetAPost(post1, username, tokenRetweeterer, r, t)
	retweet_post2 := retweetAPost(post2, username, tokenRetweeterer, r, t)

	makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{"tag5", "tag6"}, true, "", r, t)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{retweet_post2, retweet_post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_R+"&wanted_user_id="+retweeter_id, nil)
	addAuthorization(getFeed, tokenRetweeterer)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)

	for i := range result.Data {
		checkRetweetPost(result.Data[i], username, t)
	}
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestFeedRetweetNextOffset(t *testing.T) {
	log.Println("TestFeedRetweetNextOffset")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{service.TEST_TAG_ONE, "tag5"}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_ONE, "content2", []string{"tag6", service.TEST_TAG_TWO}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_ONE, "content3", []string{service.TEST_TAG_THREE, "tag6"}, false, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	retweeter_id := service.TEST_USER_THREE
	username := service.TEST_USER_THREE_USERNAME
	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, username, true)
	assert.Equal(t, err, nil, "Error should be nil")

	retweet_post1 := retweetAPost(post1, username, tokenRetweeterer, r, t)
	time.Sleep(1 * time.Second)
	retweet_post2 := retweetAPost(post2, username, tokenRetweeterer, r, t)
	time.Sleep(1 * time.Second)
	retweet_post3 := retweetAPost(post3, username, tokenRetweeterer, r, t)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{retweet_post3, retweet_post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "2"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_R+"&wanted_user_id="+retweeter_id, nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 2, result.Pagination.Limit, "Limit should be 2")
	assert.Equal(t, 2, result.Pagination.Next_Offset, "Next offset should be 2")

	result2 := models.ReturnPaginatedPosts{}
	expectedPosts2 := []models.FrontPost{retweet_post1}

	skip_2 := strconv.Itoa(result.Pagination.Next_Offset)

	getFeed2, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip_2+"&limit="+limit+"&feed_type="+FEED_TYPE_R+"&wanted_user_id="+retweeter_id, nil)
	addAuthorization(getFeed2, token)

	feedRecorder2 := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder2, getFeed2)

	err_3 := json.Unmarshal(feedRecorder2.Body.Bytes(), &result2)

	assert.Equal(t, err_3, nil)
	assert.Equal(t, http.StatusOK, feedRecorder2.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts2, result2.Data, t)
	assert.Equal(t, 2, result2.Pagination.Limit, "Limit should be 2")
	assert.Equal(t, 0, result2.Pagination.Next_Offset, "Next offset should be 0")
}
