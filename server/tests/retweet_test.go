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

func TestRetweetAPost(t *testing.T) {
	log.Println("TestRetweetAPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	retweeter_id := service.TEST_USER_TWO

	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, service.TEST_USER_TWO_USERNAME, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	retweetPost, _ := http.NewRequest("POST", "/twitsnap/retweet/"+post.Post_ID, nil)
	addAuthorization(retweetPost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, retweetPost)

	retweet_post := models.FrontPost{}

	err = json.Unmarshal(first.Body.Bytes(), &retweet_post)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusCreated, first.Code)
	assert.Equal(t, true, retweet_post.User_Retweet)
	assert.Equal(t, retweet_post.Content, post.Content)
	assert.Equal(t, retweet_post.Tags, post.Tags)
	assert.Equal(t, retweet_post.Retweet_Author, service.TEST_USER_TWO_USERNAME)

	getPostAfterRetweet, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPostAfterRetweet, tokenRetweeterer)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostAfterRetweet)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	// log.Println(result_post)
	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, true, result_post.User_Retweet)
	assert.Equal(t, result_post.Retweets, 1)

	getPostRetweeted, _ := http.NewRequest("GET", "/twitsnap/"+retweet_post.Post_ID, nil)
	addAuthorization(getPostRetweeted, tokenRetweeterer)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, getPostRetweeted)

	retweet_result := models.FrontPost{}

	err = json.Unmarshal(third.Body.Bytes(), &retweet_result)

	// log.Println(retweet_result)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, third.Code)
	assert.Equal(t, true, retweet_result.User_Retweet)
}

func TestRetweetInFeedFollowing(t *testing.T) {
	log.Println("TestRetweetInFeedFollowing")

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

	retweetPost, _ := http.NewRequest("POST", "/twitsnap/retweet/"+post1.Post_ID, nil)
	addAuthorization(retweetPost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, retweetPost)

	retweet_post := models.FrontPost{}

	err = json.Unmarshal(first.Body.Bytes(), &retweet_post)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusCreated, first.Code)
	assert.Equal(t, true, retweet_post.User_Retweet)
	assert.Equal(t, retweet_post.Content, post1.Content)
	assert.Equal(t, retweet_post.Tags, post1.Tags)
	assert.Equal(t, retweet_post.Retweet_Author, service.TEST_USER_TWO_USERNAME)

	// log.Println(retweet_post)

	expectedPosts := []models.FrontPost{retweet_post, post3, post2, post1}

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
	checkRetweetPost(result.Data[0], service.TEST_USER_TWO_USERNAME, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)

}

func TestUnRetweetAPost(t *testing.T) {
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

	_ = retweetAPost(post, username, tokenRetweeterer, r, t)

	getPostAfterRetweet, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPostAfterRetweet, tokenRetweeterer)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostAfterRetweet)

	result_post := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	// log.Println(result_post)
	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, true, result_post.User_Retweet)
	assert.Equal(t, result_post.Retweets, 1)
	assert.Equal(t, result_post.Content, post.Content)
	assert.Equal(t, result_post.Tags, post.Tags)

	unretweetPost, _ := http.NewRequest("DELETE", "/twitsnap/retweet/"+post.Post_ID, nil)
	addAuthorization(unretweetPost, tokenRetweeterer)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, unretweetPost)

	assert.Equal(t, http.StatusNoContent, third.Code)

	getPostAfterUnRetweet, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPostAfterUnRetweet, tokenRetweeterer)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPostAfterUnRetweet)

	// result_post_no_retweet := models.FrontPost{}

	err = json.Unmarshal(fourth.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")

	// log.Println(result_post)
	
	assert.Equal(t, http.StatusOK, fourth.Code)
	assert.Equal(t, false, result_post.User_Retweet)
	assert.Equal(t, result_post.Retweets, 0)
}

func TestRetweetInFeedForyou(t *testing.T) {
	log.Println("TestRetweetInFeedForyou")

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

	retweet_post := retweetAPost(post1, username, tokenRetweeterer, r, t)

	// log.Println(retweet_post)

	expectedPosts := []models.FrontPost{retweet_post, post3, post2, post1}

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
	checkRetweetPost(result.Data[0], username, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestRetweetInSingleRetweeterFeed(t *testing.T) {
	log.Println("TestRetweetInSingleRetweeterFeed")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{service.TEST_TAG_ONE, "tag2"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_ONE, "content3 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, service.TEST_USER_ONE_USERNAME, false)

	assert.Equal(t, err, nil, "Error should be nil")

	retweeter := service.TEST_USER_ONE
	username := service.TEST_USER_ONE_USERNAME
	tokenRetweeterer, err := auth.GenerateToken(retweeter, username, false)

	assert.Equal(t, err, nil, "Error should be nil")

	retweet_post := retweetAPost(post2, username, tokenRetweeterer, r, t)

	// log.Println(retweet_post)

	expectedPosts := []models.FrontPost{retweet_post, post3, post1}

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
	checkRetweetPost(result.Data[0], username, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestRetweetAndLikeOriginal(t *testing.T) {
	log.Println("TestRetweetAndLikeOriginal")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	retweeter_id := service.TEST_USER_TWO
	username := service.TEST_USER_TWO_USERNAME

	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, username, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	likePost, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Post_ID, nil)
	addAuthorization(likePost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, likePost)

	retweet_post := retweetAPost(post, username, tokenRetweeterer, r, t)


	assert.Equal(t, http.StatusNoContent, first.Code, "Status should be 204")

	getRetweet, _ := http.NewRequest("GET", "/twitsnap/"+retweet_post.Post_ID, nil)
	addAuthorization(getRetweet, tokenRetweeterer)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getRetweet)

	retweet_result := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &retweet_result)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
	assert.Equal(t, true, retweet_result.User_Retweet, "User should have retweeted")
	assert.Equal(t, 1, retweet_result.Retweets, "Retweets should be 1")
	assert.Equal(t, 1, retweet_result.Likes, "Likes should be 1")
	assert.Equal(t, retweet_result.User_Liked, true, "Original post should be liked")
}

func TestRetweetAndLikeOriginalAndRetweet(t *testing.T) {
	log.Println("TestRetweetAndLikeOriginalAndRetweet")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE
	retweeter_id := service.TEST_USER_TWO
	username := service.TEST_USER_TWO_USERNAME

	liker_id := service.TEST_USER_THREE
	tokenLiker, err := auth.GenerateToken(liker_id, service.TEST_USER_THREE_USERNAME, true)

	assert.Equal(t, err, nil, "Error should be nil")

	tokenRetweeterer, err := auth.GenerateToken(retweeter_id, username, true)
	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	retweet_post := retweetAPost(post, username, tokenRetweeterer, r, t)

	likePost, _ := http.NewRequest("POST", "/twitsnap/like/"+post.Post_ID, nil)
	addAuthorization(likePost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, likePost)

	assert.Equal(t, http.StatusNoContent, first.Code, "Status should be 204")

	likePostViaRt, _ := http.NewRequest("POST", "/twitsnap/like/"+retweet_post.Original_Post_ID, nil)
	addAuthorization(likePostViaRt, tokenLiker)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, likePostViaRt)

	assert.Equal(t, http.StatusNoContent, third.Code, "Status should be 204")

	getRetweet, _ := http.NewRequest("GET", "/twitsnap/"+retweet_post.Post_ID, nil)
	addAuthorization(getRetweet, tokenRetweeterer)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getRetweet)

	retweet_result := models.FrontPost{}

	err = json.Unmarshal(second.Body.Bytes(), &retweet_result)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
	assert.Equal(t, true, retweet_result.User_Retweet, "User should have retweeted")
	assert.Equal(t, 1, retweet_result.Retweets, "Retweets should be 1")
	assert.Equal(t, 2, retweet_result.Likes, "Likes should be 2")
	assert.Equal(t, retweet_result.User_Liked, true, "Original post should be liked")

}

func TestRetwetSingleAuthorField(t *testing.T) {
	log.Println("TestRetwetSingleAuthorField")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	tags := []string{"tag1", "tag2"}

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_ONE, "content3 " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	token, err := auth.GenerateToken(service.TEST_USER_ONE, service.TEST_USER_ONE_USERNAME, false)

	assert.Equal(t, err, nil, "Error should be nil")

	retweeter := service.TEST_USER_THREE
	username := service.TEST_USER_THREE_USERNAME
	tokenRetweeterer, err := auth.GenerateToken(retweeter, username, false)

	assert.Equal(t, err, nil, "Error should be nil")

	_ = retweetAPost(post2, username, tokenRetweeterer, r, t)

	// log.Println(retweet_post)

	expectedPosts := []models.FrontPost{post3, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+FEED_TYPE_S+"&wanted_user_id="+service.TEST_USER_ONE, nil)
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