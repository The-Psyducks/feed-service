package test

import (
	"log"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"server/src/auth"
	"server/src/models"
	"server/src/router"
)

func TestNewPost(t *testing.T) {

	log.Println("TestNewPost")

	db := ConnectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "content", Tags: []string{"tag1", "tag2"}, Public: true}
	req := NewPostRequest(postBody, r)

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	AddAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := models.FrontPost{}

	err = json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil)
	MakeResponseAsserions(t, http.StatusCreated, result, postBody, author_id, first.Code)
}