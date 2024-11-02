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
	FEED_TYPE_Y = "foryou"
)

func TestFeedForYou(t *testing.T) {
	log.Println("TestGetFeedForyou")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{service.TEST_TAG_ONE, "tag5"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags2 := []string{"tag6", service.TEST_TAG_TWO}

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + "#" + tags2[0] + " #" + tags2[1], tags2, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags3 := []string{service.TEST_TAG_THREE, "tag6"}

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3 " + "#" + tags3[0] + " #" + tags3[1], tags3, true, "", r, t)

	tags4 := []string{"tag7", "tag8"}

	makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content4 " + "#" + tags4[0] + " #" + tags4[1], tags4, true, "", r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post3, post2, post1}

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

	log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestFeedForyouNotFollowing(t *testing.T) {
	log.Println("TestGetFeedForyouNotFollowing")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{service.TEST_TAG_ONE, "tag5"}

	post1 := makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags2 := []string{"tag6", service.TEST_TAG_TWO}

	post2 := makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content " + "#" + tags2[0] + " #" + tags2[1], tags2, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags3 := []string{service.TEST_TAG_THREE, "tag6"}

	_ = makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content " + "#" + tags3[0] + " #" + tags3[1], tags3, false, "", r, t)

	tags4 := []string{"tag7", "tag8"}

	makeAndAssertPost(service.TEST_USER_THREE, "content " + "#" + tags4[0] + " #" + tags4[1], tags4, true, "", r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post2, post1}

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
	assertOnlyPublicPosts(result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestFeedForyouNextOffset(t *testing.T) {
	log.Println("TestGetFeedFollowingNextOffset")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{service.TEST_TAG_ONE, "tag5"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags2 := []string{"tag6", service.TEST_TAG_TWO}

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content " + "#" + tags2[0] + " #" + tags2[1], tags2, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags3 := []string{service.TEST_TAG_THREE, "tag6"}

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content " + "#" + tags3[0] + " #" + tags3[1], tags3, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post3, post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "2"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_Y+"&wanted_user_id="+"", nil)
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

	getFeed2, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip_2+"&limit="+limit+"&feed_type="+FEED_TYPE_Y+"&wanted_user_id="+"", nil)
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
