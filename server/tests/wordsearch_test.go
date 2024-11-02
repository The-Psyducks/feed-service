package test

import (
	"log"
	"strconv"
	"strings"
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

func TestWordsearch(t *testing.T) {

	log.Println("TestWordsearch: post with one or more of the words wanted in the content appear")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	words_wanted := "apple pie pecan"

	words_wanted_list := strings.Split(words_wanted, " ")

	tags := []string{service.TEST_TAG_ONE, "tag5"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[0], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[1], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3 " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[2], tags, true, "", r, t)

	
	makeAndAssertPost(service.TEST_USER_THREE, "content4 " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post3, post2, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/word-search?words="+words_wanted+"&time="+time+"&skip="+skip+"&limit="+limit, nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)
	postsHaveAtLeastOneWord(result.Data, words_wanted_list, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestWordSearchNotFollowing(t *testing.T) {

	log.Println("TestHashagSearchNotFollowing: only posts with all the tags wanted appear, and only public posts from not followed users")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	words_wanted := "apple pie pecan"

	words_wanted_list := strings.Split(words_wanted, " ")

	tags := []string{"tags_wanted", "tags"}

	makeAndAssertPost(service.TEST_NOT_FOLLOWING_ID, "content " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[2], tags, false, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[0], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{}, true, "", r, t)

	makeAndAssertPost(service.TEST_USER_THREE, "content4", []string{}, true, "", r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/word-search?words="+words_wanted+"&time="+time+"&skip="+skip+"&limit="+limit, nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assertOnlyPublicPostsForNotFollowing(result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestWordSearchFollowing(t *testing.T) {

	log.Println("TestHashtagSearchFollowing: only posts with all the tags wanted appear, with public and not public posts from followed users")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	words_wanted := "apple pie pecan"

	words_wanted_list := strings.Split(words_wanted, " ")

	tags := []string{"tags_wanted", "tags_wanted[1]"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[1], tags, false, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + "#" + tags[0] + " #" + tags[1] + " " + words_wanted_list[0], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	_ = makeAndAssertPost(service.TEST_USER_THREE, "content " + "#" + tags[0] + " #" + tags[1] + " ", tags, true, "", r, t)

	makeAndAssertPost(service.TEST_USER_THREE, "content4", []string{}, true, "", r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post2, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/word-search?words="+words_wanted+"&time="+time+"&skip="+skip+"&limit="+limit, nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil, "Error should be nil")
	assert.Equal(t, http.StatusOK, feedRecorder.Code, "Status should be 200")

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assertOnlyPublicPostsForNotFollowing(result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestWordSearchNextOffset(t *testing.T) {
	log.Println("TestGetFeedFollowingNextOffset")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	words_wanted := "apple pie pecan"

	words_wanted_list := strings.Split(words_wanted, " ")

	tags := []string{service.TEST_TAG_ONE, "tag5"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[0], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content22 " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[1], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3 " + "#" + tags[0] + " #" + tags[1] + " " +words_wanted_list[2], tags, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	expectedPosts := []models.FrontPost{post3, post2}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "2"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/word-search?words="+words_wanted+"&time="+time+"&skip="+skip+"&limit="+limit, nil)
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

	getFeed2, _ := http.NewRequest("GET", "/twitsnap/word-search?words="+words_wanted+"&time="+time+"&skip="+skip_2+"&limit="+limit, nil)
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
