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

	postErrors "server/src/all_errors"
)

const (
	FEED_TYPE_F = "following"
)

func TestFeedFollowing(t *testing.T) {
	log.Println("TestGetFeedFollowing")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{"tag1", "tag2"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags2 := []string{"tag3", "tag4"}

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content 2 " + "#" + tags2[0] + " #" + tags2[1], tags2, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags3 := []string{"tag5", "tag6"}

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content 3 " + "#" + tags3[0] + " #" + tags3[1], tags3, []string{}, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post3, post2, post1}

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
	assert.Equal(t, 6, result.Pagination.Limit, "Limit should be 6")
	assert.Equal(t, 0, result.Pagination.Next_Offset, "Next offset should be 0")
}

func TestFeedFollowingNextOffset(t *testing.T) {
	log.Println("TestGetFeedFollowingNextOffset")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{"tag1", "tag2"}

	post1 := makeAndAssertPost(service.TEST_USER_TWO, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags2 := []string{"tag3", "tag4"}

	post2 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags2[0] + " #" + tags2[1], tags2, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags3 := []string{"tag5", "tag6"}

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content " + "#" + tags3[0] + " #" + tags3[1], tags3, []string{}, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post3, post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "2"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_F+"&wanted_user_id="+"", nil)
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
	expectedPosts2 := []models.FrontPost{post1}

	skip_2 := strconv.Itoa(result.Pagination.Next_Offset)

	getFeed2, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip_2+"&limit="+limit+"&feed_type="+FEED_TYPE_F+"&wanted_user_id="+"", nil)
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

func TestFeedBadRequestFollowing(t *testing.T) {
	log.Println("FeedBadRequestFollowing")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, "username", false)

	assert.Equal(t, err, nil, "Error should be nil")
	
	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)
	
	skip := "0"
	limit := "6"
	feed := "following_bad"
	
	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed, nil)
	addAuthorization(getFeed, token)
	
	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)
	
	result := postErrors.BadFeedRequest(feed)
	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusBadRequest, feedRecorder.Code, "Status should be 400")
	assert.Equal(t, postErrors.BadFeedRequest(feed), result, "Error should be the same")
}