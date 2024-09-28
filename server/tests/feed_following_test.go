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

func TestGetFeedFollowing(t *testing.T) {
	log.Println("TestGetFeedFollowing")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{"tag1", "tag2"}, true, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{"tag3", "tag4"}, true, r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{"tag5", "tag6"}, true, r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post3, post2, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "following"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 6, result.Limit)
	assert.Equal(t, 0, result.Next_Offset)
}

func TestGetFeedFollowingNextOffset(t *testing.T) {
	log.Println("TestGetFeedFollowingNextOffset")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_TWO, "content", []string{"tag1", "tag2"}, true, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_ONE, "content2", []string{"tag3", "tag4"}, true, r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{"tag5", "tag6"}, true, r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post3, post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "following"
	limit := "2"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 2, result.Limit)
	assert.Equal(t, 2, result.Next_Offset)

	result2 := models.ReturnPaginatedPosts{}
	expectedPosts2 := []models.FrontPost{post1}

	skip_2 := strconv.Itoa(result.Next_Offset)

	getFeed2, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip_2+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed2, token)

	feedRecorder2 := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder2, getFeed2)

	err_3 := json.Unmarshal(feedRecorder2.Body.Bytes(), &result2)

	assert.Equal(t, err_3, nil)
	assert.Equal(t, http.StatusOK, feedRecorder2.Code)

	compareOrderAsExpected(expectedPosts2, result2.Data, t)
	assert.Equal(t, 2, result2.Limit)
	assert.Equal(t, 0, result2.Next_Offset)
}