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

func TestGetAll(t *testing.T) {
	log.Println("TestGetAll")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{service.TEST_TAG_ONE, "tag5"}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{"tag6", service.TEST_TAG_TWO}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{service.TEST_TAG_THREE, "tag6"}, true, "", r, t)

	post4 := makeAndAssertPost(service.TEST_USER_THREE, "content4", []string{"tag7", "tag8"}, true, "", r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post4, post3, post2, post1}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getAll, _ := http.NewRequest("GET", "/twitsnap/all?time="+time+"&skip="+skip+"&limit="+limit, nil)
	addAuthorization(getAll, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getAll)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 6, result.Pagination.Limit)
	assert.Equal(t, 0, result.Pagination.Next_Offset)
}

func TestGetAllDeniedAccess(t *testing.T) {
	log.Println("TestGetAllDeniedAccess")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	token, err := auth.GenerateToken("1", "username", false)

	assert.Equal(t, err, nil)

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "6"

	getAll, _ := http.NewRequest("GET", "/twitsnap/all?time="+time+"&skip="+skip+"&limit="+limit, nil)
	addAuthorization(getAll, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getAll)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusForbidden, feedRecorder.Code)
}

func TestGetAllNextOffset(t *testing.T) {
	log.Println("TestGetAllNextOffset")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	post1 := makeAndAssertPost(service.TEST_USER_ONE, "content", []string{service.TEST_TAG_ONE, "tag5"}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post2 := makeAndAssertPost(service.TEST_USER_TWO, "content2", []string{"tag6", service.TEST_TAG_TWO}, true, "", r, t)

	time.Sleep(1 * time.Second)

	post3 := makeAndAssertPost(service.TEST_USER_THREE, "content3", []string{service.TEST_TAG_THREE, "tag6"}, true,  "",r, t)

	post4 := makeAndAssertPost(service.TEST_USER_THREE, "content4", []string{"tag7", "tag8"}, true, "", r, t)

	token, err := auth.GenerateToken("1", "username", true)

	assert.Equal(t, err, nil)

	expectedPosts := []models.FrontPost{post4, post3}

	result := models.ReturnPaginatedPosts{}

	time.Sleep(1 * time.Second)
	time := time.Now().Format(time.RFC3339)

	skip := "0"
	limit := "2"

	getAll, _ := http.NewRequest("GET", "/twitsnap/all?time="+time+"&skip="+skip+"&limit="+limit, nil)
	addAuthorization(getAll, token)

	feedRecorder := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder, getAll)

	err_2 := json.Unmarshal(feedRecorder.Body.Bytes(), &result)

	// log.Println(result)

	assert.Equal(t, err_2, nil)
	assert.Equal(t, http.StatusOK, feedRecorder.Code)

	compareOrderAsExpected(expectedPosts, result.Data, t)
	assert.Equal(t, 2, result.Pagination.Limit)
	assert.Equal(t, 2, result.Pagination.Next_Offset)

	skip_2 := strconv.Itoa(result.Pagination.Next_Offset)

	expectedPosts2 := []models.FrontPost{post2, post1}

	result2 := models.ReturnPaginatedPosts{}

	getAll2, _ := http.NewRequest("GET", "/twitsnap/all?time="+time+"&skip="+skip_2+"&limit="+limit, nil)
	addAuthorization(getAll2, token)

	feedRecorder2 := httptest.NewRecorder()
	r.ServeHTTP(feedRecorder2, getAll2)

	err_3 := json.Unmarshal(feedRecorder2.Body.Bytes(), &result2)

	assert.Equal(t, err_3, nil)
	assert.Equal(t, http.StatusOK, feedRecorder2.Code)

	compareOrderAsExpected(expectedPosts2, result2.Data, t)
	assert.Equal(t, 2, result2.Pagination.Limit)
	assert.Equal(t, 0, result2.Pagination.Next_Offset)
}