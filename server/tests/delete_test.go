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

	assert.Equal(t, err, nil)

	post := makeAndAssertPost("1", "content", []string{"tag1", "tag2"}, true, r, t)

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(deletePost, token)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNoContent, third.Code)

	getPost, _ := http.NewRequest("GET", "/twitsnap/"+post.Post_ID, nil)
	addAuthorization(getPost, token)

	fourth := httptest.NewRecorder()
	r.ServeHTTP(fourth, getPost)
	assert.Equal(t, http.StatusNotFound, fourth.Code)
}

func TestDeleteUnexistentPost(t *testing.T) {

	log.Println("TestDeleteUnexistentPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	token, err := auth.GenerateToken("1", "username", false)

	assert.Equal(t, err, nil)

	post := makeAndAssertPost("1", "content", []string{"tag1", "tag2"}, true, r, t)

	deletePost, _ := http.NewRequest("DELETE", "/twitsnap/"+post.Post_ID+"invalid", nil)
	addAuthorization(deletePost, token)

	third := httptest.NewRecorder()
	r.ServeHTTP(third, deletePost)

	assert.Equal(t, http.StatusNotFound, third.Code)
}
