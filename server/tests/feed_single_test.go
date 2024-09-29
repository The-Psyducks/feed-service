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

func TestGetFeedSingle(t *testing.T) {
	log.Println("TestGetFeedSingle")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{"tag1", "tag2"}, true, r, t)

	time.Sleep(1 * time.Second)

	makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{"tag3", "tag4"}, true, r, t)

	time.Sleep(1 * time.Second)

	makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{"tag5", "tag6"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "single"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+"1", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestGetFeedSingleNotFollowing(t *testing.T) {
	log.Println("TestGetFeedSingleNotFollowing: should only get the public posts of the user")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	authorId := service.TEST_NOT_FOLLOWING_ID

	post1 := makeAndAssertPost(authorId, "content", []string{"tag1", "tag2"}, true, r, t)

	time.Sleep(1 * time.Second)

	makeAndAssertPost(authorId, "content2", []string{"tag3", "tag4"}, false, r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(authorId, "content3", []string{"tag5", "tag6"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post3, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "single"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+authorId, nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assertOnlyPublicPosts(result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestFeedSingleWithFollowing(t *testing.T) {
	log.Println("TestGetFeedSingleNotFollowing: should only get the public posts of the user not folowing and all the posts of the one followed")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	authorId := service.TEST_USER_ONE

	post1 := makeAndAssertPost(authorId, "content", []string{"tag1", "tag2"}, true, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(authorId, "content2", []string{"tag3", "tag4"}, false, r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(authorId, "content3", []string{"tag5", "tag6"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post3, post2, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "single"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+authorId, nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}