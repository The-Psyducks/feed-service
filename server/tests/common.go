package test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/src/database"
	"server/src/models"
	"testing"

	// "server/src/router"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostBody struct {
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Public  bool     `json:"public"`
}

func NewPostRequest(post PostBody, r *gin.Engine) *http.Request {
	marshalledData, _ := json.Marshal(post)
	req, _ := http.NewRequest("POST", "/twitsnap", bytes.NewReader(marshalledData))

	req.Header.Add("content-type", "application/json")

	return req
}

func AddAuthorization(req *http.Request, token string) {
	req.Header.Add("Authorization", "Bearer "+token)
}

func ConnectToDatabase() database.Database {

	log.Println("Connect to database")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	db := database.NewAppDatabase(client)

	return db
}

func MakeResponseAsserions(t *testing.T, response int, result_post models.FrontPost, postBody PostBody, author_id string, code int) {
	assert.Equal(t, response, code)
	assert.Equal(t, result_post.Content, postBody.Content)
	assert.Equal(t, result_post.Author_Info.Author_ID, author_id)
	assert.Equal(t, result_post.Tags, postBody.Tags)
	assert.Equal(t, result_post.Public, postBody.Public)
}
