package test

import (
	"log"
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

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content " + words_wanted_list[0], []string{service.TEST_TAG_ONE, "tag5"}, true, r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2 " + words_wanted_list[1], []string{"tag6", "tag5"}, true, r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3 " + words_wanted_list[2], []string{service.TEST_TAG_THREE, "tag6"}, true, r, t)

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

	getFeed, _ := http.NewRequest("GET", "/twitsnap/wordsearch?words="+words_wanted+"&time="+time+"&skip="+skip+"&limit="+limit+"&feed_type="+feed_type+"&wanted_user_id="+"", nil)
	addAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	postsHaveAtLeastOneWord(result.Data, words_wanted_list, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}