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

	// "time"

	"github.com/stretchr/testify/assert"
)

func TestBookmarkPost(t *testing.T) {
	log.Println("TestBookmarkPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	bookmarker_id := service.TEST_USER_TWO

	tokenLiker, err := auth.GenerateToken(bookmarker_id, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	bookmarkPost, _ := http.NewRequest("POST", "/twitsnap/bookmark/"+post.Original_Post_ID, nil)
	addAuthorization(bookmarkPost, tokenLiker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, bookmarkPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	getPostLiked, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getPostLiked, tokenLiker)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostLiked)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
	assert.Equal(t, result_post.Bookmark, true, "Post should be bookmarked")
}

func TestUnBookmarkPost(t *testing.T) {
	log.Println("TestUnBookmarkPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	bookmarker_id := service.TEST_USER_TWO

	tokenLiker, err := auth.GenerateToken(bookmarker_id, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	bookmarkPost, _ := http.NewRequest("POST", "/twitsnap/bookmark/"+post.Original_Post_ID, nil)
	addAuthorization(bookmarkPost, tokenLiker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, bookmarkPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	getBookmarkedPost, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getBookmarkedPost, tokenLiker)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getBookmarkedPost)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
	assert.Equal(t, result_post.Bookmark, true, "Post should be bookmarked")


	getPost2, _ := http.NewRequest("DELETE", "/twitsnap/bookmark/"+post.Original_Post_ID, nil)
	addAuthorization(getPost2, tokenLiker)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getPost2)

	assert.Equal(t, http.StatusNoContent, third.Code)

	getUnBookmarkedPost, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getUnBookmarkedPost, tokenLiker)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getUnBookmarkedPost)

	result_post_no_bookmark := models.FrontPost{}

	err = json.Unmarshal(fourth.Body.Bytes(), &result_post_no_bookmark)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, fourth.Code, "Status should be 200")
	assert.Equal(t, result_post_no_bookmark.Bookmark, false, "Post should not be bookmarked")
}

func TestFetchBookmarkedPosts(t *testing.T) {
	log.Println("TestFetchBookmarkedPosts")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	bookmarker_id := service.TEST_USER_TWO

	tokenLiker, err := auth.GenerateToken(bookmarker_id, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content 3 " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	boookmarkPost(post1.Original_Post_ID, tokenLiker, r, t)

	boookmarkPost(post2.Original_Post_ID, tokenLiker, r, t)

	boookmarkPost(post3.Original_Post_ID, tokenLiker, r, t)

	time.Sleep(1 * time.Second)

	time := time.Now().Format(time.RFC3339)
	skip := "0"
	limit := "6"

	expectedPosts := []models.FrontPost{post3, post2, post1}

	getBookmarkedPosts, _ := http.NewRequest("GET", "/twitsnap/bookmarks?time="+time+"&skip="+skip+"&limit="+limit+"&wanted_user_id="+ bookmarker_id, nil)

	addAuthorization(getBookmarkedPosts, tokenLiker)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getBookmarkedPosts)

	result := models.ReturnPaginatedPosts{}

	err = json.Unmarshal(second.Body.Bytes(), &result)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
	compareOrderAsExpected(expectedPosts, result.Data, t)
	
	assert.Equal(t, result.Data[0].Bookmark, true, "Post should be bookmarked")
	assert.Equal(t, result.Data[1].Bookmark, true, "Post should be bookmarked")
	assert.Equal(t, result.Data[2].Bookmark, true, "Post should be bookmarked")

	assert.Equal(t, 6, result.Pagination.Limit, "Limit should be 6")
	assert.Equal(t, 0, result.Pagination.Next_Offset, "Next offset should be 0")
}

func TestBookmarkReTweetedPost(t *testing.T) {
	log.Println("TestBookmarkPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	bookmarker_id := service.TEST_USER_TWO

	tokenBookmarker, err := auth.GenerateToken(bookmarker_id, service.TEST_USER_TWO_USERNAME, true)

	assert.Equal(t, err, nil, "Error should be nil")

	tokenRetweeter, err := auth.GenerateToken(service.TEST_USER_THREE, service.TEST_USER_THREE_USERNAME, true)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, true, "", r, t)

	reTweetPost := retweetAPost(post, service.TEST_USER_THREE_USERNAME, tokenRetweeter, r, t)

	bookmarkPost, _ := http.NewRequest("POST", "/twitsnap/bookmark/"+reTweetPost.Original_Post_ID, nil)
	addAuthorization(bookmarkPost, tokenBookmarker)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, bookmarkPost)

	assert.Equal(t, http.StatusNoContent, first.Code)

	getPostBookmarked, _ := http.NewRequest("GET", "/twitsnap/"+post.Original_Post_ID, nil)
	addAuthorization(getPostBookmarked, tokenBookmarker)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostBookmarked)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
	assert.Equal(t, result_post.Bookmark, true, "Post should be bookmarked")

	getPostRetweet, _ := http.NewRequest("GET", "/twitsnap/"+reTweetPost.Post_ID, nil)
	addAuthorization(getPostRetweet, tokenBookmarker)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getPostRetweet)

	result_rt_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_rt_post)

	// log.Println(result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, third.Code, "Status should be 200")
	assert.Equal(t, result_rt_post.Bookmark, true, "Post should be bookmarked")
}