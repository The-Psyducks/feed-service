package test

import (
	"testing"

	// "bytes"
	// "encoding/json"
	// "net/http"
	// "net/http/httptest"

	// // "github.com/gin-gonic/gin"
	// "github.com/mjarkk/mongomock"
	"github.com/stretchr/testify/assert"

	// "server/src/post"
	// "server/src/router"
    // "server/src/database"
)

func TestDummy(t *testing.T) {
    assert.Equal(t, 1, 1)
}

func TestNewPost(t *testing.T) {
    // create a new mock database
	// mockClient := mongomock.NewDB()

    // db := database.NewDatabase(mockClient)

    // r := router.CreateRouter(db)

    // // create a new post
    // post := post.PostExpectedFromat{
    //     Author_ID: "1",
    //     Content: "This is a twitsnap",
    //     Tags: []string{"#twitsnap", "#golang"},
    // }

    // marshalledData, _ := json.Marshal(post)
    // req, _ := http.NewRequest("POST", "/twitsnap", bytes.NewBuffer(marshalledData))

    // req.Header.Add("content-type", "application/json")

    // first := httptest.NewRecorder()
	// r.ServeHTTP(first, req)

    // result := struct {
	// 	Post_Id string `json:"post_id"`
	// 	Content  string  `json:"content"`
    //     Tags []string `json:"tags"`
    //     Author_ID string `json:"author_id"`
	// }{}

	// err := json.Unmarshal(first.Body.Bytes(), &result)

	// assert.Equal(t, err, nil)
	// assert.Equal(t, http.StatusCreated, first.Code)
	// assert.Equal(t, result.Content, post.Content)

    
}