package test

import (
	"log"
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
func TestFeedForYou(t *testing.T) {
	log.Println("TestGetFeedForyou")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{service.TEST_TAG_ONE, "tag5"}, true, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{"tag6", service.TEST_TAG_TWO}, true, r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{service.TEST_TAG_THREE, "tag6"}, true, r, t)

	makeAndAssertPost(service.TEST_USER_THREE, "content4", []string{"tag7", "tag8"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post3, post2, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "foryou"
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


func TestFeedForyouNotFollowing(t *testing.T) {
	log.Println("TestGetFeedForyouNotFollowing")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content", []string{service.TEST_TAG_ONE, "tag5"}, true, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content2", []string{"tag6", service.TEST_TAG_TWO}, true, r, t)

	time.Sleep(1 * time.Second)

	_ = makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content3", []string{service.TEST_TAG_THREE, "tag6"}, false, r, t)

	makeAndAssertPost(service.TEST_USER_THREE, "content4", []string{"tag7", "tag8"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post2, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "foryou"
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
	assertOnlyPublicPosts(result.Data, t)
	assert.Equal(t, 6, result.Limit)
	assert.Equal(t, 0, result.Next_Offset)
}