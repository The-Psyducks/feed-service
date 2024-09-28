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
)


func TestGetFeedFollowing(t *testing.T) {
	log.Println("TestGetFeedFollowing")

	db := ConnectToDatabase()

	r := router.CreateRouter(db)
	
	post1 := MakeAndAssertPost("2", "content", []string{"tag1", "tag2"}, true, r, t)

	time.Sleep(1 * time.Second)

	post2 := MakeAndAssertPost("1", "content2", []string{"tag3", "tag4"}, true, r, t)

	time.Sleep(1 * time.Second)

	post3 := MakeAndAssertPost("3", "content3", []string{"tag5", "tag6"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post3, post2, post1}

    result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	feed_type := "following"
	limit := "6"
	
	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time=" + time + "&skip="+skip+"&limit="+limit+"&feed_type="+ feed_type + "&wanted_user_id="+"", nil)
	AddAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)
	
	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	log.Println(result)
	
	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	
	assert.Equal(t, true, compareOrderAsExpected(expectedPosts, result.Data))
	assert.Equal(t, 6, result.Limit)
	assert.Equal(t, 0, result.Next_Offset)
}

func compareOrderAsExpected(expected []models.FrontPost, result []models.FrontPost) bool {
	if len(expected) != len(result) {
		return false
	}
	for i := range expected {
		if expected[i].Content != result[i].Content {
			return false
		}
	}
	return true
}

func TestGetFeedSingle(t *testing.T) {
	log.Println("TestGetFeedSingle")

	db := ConnectToDatabase()

	r := router.CreateRouter(db)
	
	post1 := MakeAndAssertPost("1", "content", []string{"tag1", "tag2"}, true, r, t)

	time.Sleep(1 * time.Second)

	MakeAndAssertPost("2", "content2", []string{"tag3", "tag4"}, true, r, t)

	time.Sleep(1 * time.Second)

	MakeAndAssertPost("3", "content3", []string{"tag5", "tag6"}, true, r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post1}

    result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)
	
	skip := "0"
	feed_type := "single"
	limit := "6"
	
	getFeed, _ := http.NewRequest("GET", "/twitsnap/feed?time=" + time + "&skip="+skip+"&limit="+limit+"&feed_type="+ feed_type + "&wanted_user_id="+"1", nil)
	AddAuthorization(getFeed, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getFeed)
	
	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	log.Println(result)
	
	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	
	assert.Equal(t, true, compareOrderAsExpected(expectedPosts, result.Data))
	assert.Equal(t, 6, result.Limit)
	assert.Equal(t, 0, result.Next_Offset)
}