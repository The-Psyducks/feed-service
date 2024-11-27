package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"server/src/auth"
	"server/src/router"
	"server/src/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTrending(t *testing.T) {
	log.Println("TestGetTrending")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := service.TEST_USER_ONE

	tokenAsker, err := auth.GenerateToken(author_id, "username", true)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	makeAndAssertPost(author_id, "content " + "#" + tags[0] + " #" + tags[1], tags, []string{}, true, "", r, t)

	time.Sleep(1 * time.Second)

	tags2 := []string{tags[0]}

	makeAndAssertPost(author_id, "content " + "#" + tags2[0], tags2, []string{}, true, "", r, t)

	getPostLiked, _ := http.NewRequest("GET", "/twitsnap/trending", nil)
	addAuthorization(getPostLiked, tokenAsker)

	second := httptest.NewRecorder()
	r.ServeHTTP(second, getPostLiked)

	result_post := []string{}

	err = json.Unmarshal(second.Body.Bytes(), &result_post)

	assert.Equal(t, err, nil, "Error should be nil")
	
	assert.Equal(t, http.StatusOK, second.Code, "Status should be 200")
	assert.Equal(t, len(result_post), 2, "Should have 2 trending topics")
	
	for i := 0; i < len(tags); i++ {
		assert.Equal(t, result_post[i], tags[i], "Trending topic should be " + tags[i])
	}
}