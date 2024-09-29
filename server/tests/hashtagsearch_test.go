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

func TestHashagSearch(t *testing.T) {

	log.Println("TestHashagSearch: only posts with all the tags wanted appear")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	_ = makeAndAssertPost(service.TEST_USER_ONE, "content", []string{service.TEST_TAG_ONE, "tag5"}, true, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{"tag6", "tag5"}, true, r, t)

	time.Sleep(1 * time.Second)

	_ = makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{service.TEST_TAG_THREE, "tag6"}, true, r, t)

	makeAndAssertPost(service.TEST_USER_THREE, "content4", []string{"tag7", "tag8"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "foryou"
	limit := "6"

	tags_wanted := []string{"tag5", "tag6"}

	getFeed, _ := http.NewRequest("GET", "/twitsnap/hashtags?tags="+tags_wanted[0]+"&tags=" + tags_wanted[1] +"&time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assertOnlyPublicPostsForNotFollowing(result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestHashagSearchNotFollowing(t *testing.T) {

	log.Println("TestHashagSearchNotFollowing: only posts with all the tags wanted appear, and only public posts from not followed users")

	db := connectToDatabase()
	
	r := router.CreateRouter(db)
	
	tags_wanted := []string{"tag5", "tag6"}

	makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content", []string{tags_wanted[0], tags_wanted[1]}, false, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{tags_wanted[0], tags_wanted[1]}, true, r, t)

	time.Sleep(1 * time.Second)

	makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{service.TEST_TAG_THREE, "tag6"}, true, r, t)

	makeAndAssertPost(service.TEST_USER_THREE, "content4", []string{"tag7", "tag8"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "foryou"
	limit := "6"


	getFeed, _ := http.NewRequest("GET", "/twitsnap/hashtags?tags="+tags_wanted[0]+"&tags=" + tags_wanted[1] +"&time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assertOnlyPublicPostsForNotFollowing(result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestHashtagSearchFollowing(t *testing.T) {

	log.Println("TestHashtagSearchFollowing: only posts with all the tags wanted appear, with public and not public posts from followed users")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags_wanted := []string{"tag5", "tag6"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{tags_wanted[0], tags_wanted[1]}, false, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{tags_wanted[0], tags_wanted[1]}, true, r, t)

	time.Sleep(1 * time.Second)

	_ = makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{service.TEST_TAG_THREE, "tag6"}, true, r, t)

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


	getFeed, _ := http.NewRequest("GET", "/twitsnap/hashtags?tags="+tags_wanted[0]+"&tags=" + tags_wanted[1] +"&time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assertOnlyPublicPostsForNotFollowing(result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

