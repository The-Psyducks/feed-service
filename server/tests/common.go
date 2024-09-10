package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func NewPostRequest(author_id string, content string, tags []string, public bool, r *gin.Engine) *httptest.ResponseRecorder {
	post := struct {
		Author_ID string `json:"author_id"`
		Content  string  `json:"content"`
		Tags []string `json:"tags"`
		Public bool `json:"public"`
	}{
        Author_ID: author_id,
        Content: content,
        Tags: tags,
		Public: public,
    }
	
	marshalledData, _ := json.Marshal(post)
	req, _ := http.NewRequest("POST", "/twitsnap", bytes.NewReader(marshalledData))

	req.Header.Add("content-type", "application/json")


	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	return first
}