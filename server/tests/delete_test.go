package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"server/src/auth"
	"server/src/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletePost(t *testing.T) {

	log.Println("TestDeletePost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	token, err := auth.GenerateToken("1", "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	post := makeAndAssertPost("1", "content #tag1 #tag2", []string{"tag1", "tag2"}, []string{}, true, "", r, t)

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(deletePost, token)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNoContent, third.Code, "Status should be 204")

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPost, token)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPost)
	assert.Equal(t, http.StatusNotFound, fourth.Code, "Status should be 404")
}

func TestDeleteUnexistentPost(t *testing.T) {

	log.Println("TestDeleteUnexistentPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	token, err := auth.GenerateToken("1", "username", false)

	assert.Equal(t, err, nil, "Error should be nil")

	tags := []string{"tag1", "tag2"}

	post := makeAndAssertPost("1", "content " + "#" + tags[0] + " #" + tags[1], []string{"tag1", "tag2"}, []string{}, true, "", r, t)

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+post.Post_ID+"invalid", nil)
	addAuthorization(deletePost, token)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNotFound, third.Code, "Status should be 404")
}
