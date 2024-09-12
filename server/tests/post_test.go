package test

import (
	"testing"
	"time"

	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"

	"server/src/database"
	"server/src/router"
)

func TestNewPost(t *testing.T) {

    db := database.NewTestDatabase()

    r := router.CreateRouter(db)

	author_id := "1"
	content := "content"
	tags := []string{"tag1", "tag2"}
	public := true

    first := NewPostRequest(author_id, content,tags,public, r)

    result := struct {
		Post struct {
			Post_ID   string    `bson:"post_id"`
			Content   string    `bson:"content"`
			Author_ID string    `bson:"author_id"`
			Time      time.Time `bson:"time"`
			Public   bool    `bson:"public"`
			Tags     []string  `bson:"tags"`
		}
	}{}

	err := json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusCreated, first.Code)
	assert.Equal(t, result.Post.Content, content)
	assert.Equal(t, result.Post.Author_ID, author_id)
	assert.Equal(t, result.Post.Tags, tags)
	assert.Equal(t, result.Post.Public, public)
}