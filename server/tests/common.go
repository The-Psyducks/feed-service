package test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"server/src/auth"
	"server/src/database"
	"server/src/models"
	"server/src/service"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostBody struct {
	Content  string   `json:"content"`
	Tags     []string `json:"-"`
	Public   bool     `json:"public"`
	MediaURL string   `json:"media_url"`
	Mentions []string `json:"mentions"`
}

func newPostRequest(post PostBody) *http.Request {
	marshalledData, _ := json.Marshal(post)
	req, _ := http.NewRequest("POST", "/twitsnap", bytes.NewReader(marshalledData))

	req.Header.Add("content-type", "application/json")

	return req
}

func addAuthorization(req *http.Request, token string) {
	req.Header.Add("Authorization", "Bearer "+token)
}

func connectToDatabase() database.Database {

	log.Println("Connect to database")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	db := database.NewAppDatabase(client)

	err_2 := db.ClearDB()

	if err_2 != nil {
		log.Fatal("Error clearing database: ", err)
	}

	return db
}

func makeResponseAsserions(t *testing.T, response int, result_post models.FrontPost, postBody PostBody, author_id string, code int) {

	assert.Equal(t, response, code, "Response should be 201")
	assert.Equal(t, result_post.Content, postBody.Content, "Content should be the same")
	assert.Equal(t, result_post.Author_Info.Author_ID, author_id, "Author should be the same")
	assert.Equal(t, result_post.Tags, postBody.Tags, "Tags should be the same")
	assert.Equal(t, result_post.Public, postBody.Public, "Public should be the same")
	assert.Equal(t, result_post.Media_URL, postBody.MediaURL, "Media URL should be the same")
	assert.Equal(t, result_post.Mentions, postBody.Mentions, "Mentions should be the same")
}

func makeAndAssertPost(authorId string, content string, tags []string, mentions []string, public bool, media_url string, r *gin.Engine, t *testing.T) models.FrontPost {

	postBody := PostBody{Content: content, Tags: tags, Public: public, MediaURL: media_url, Mentions: mentions}
	req := newPostRequest(postBody)

	token, err := auth.GenerateToken(authorId, "username", false)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	addAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := models.FrontPost{}

	err = json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")
	makeResponseAsserions(t, http.StatusCreated, result, postBody, authorId, first.Code)

	return result
}

func compareOrderAsExpected(expected []models.FrontPost, result []models.FrontPost, t *testing.T) {
	assert.Equal(t, len(expected), len(result), "Length should be the same")
	for i := range expected {
		assert.Equal(t, expected[i].Content, result[i].Content, "Content should be the same")
	}
}

func assertOnlyPublicPosts(result []models.FrontPost, t *testing.T) {
	for i := range result {
		assert.Equal(t, true, result[i].Public, "Posts should be public")
	}
}

func assertOnlyPublicPostsForNotFollowing(result []models.FrontPost, t *testing.T) {
	for i := range result {
		if result[i].Author_Info.Author_ID == service.TEST_NOT_FOLLOWING_ID {
			assert.Equal(t, true, result[i].Public, "All posts should be public")
		}
	}
}

func postsHaveAtLeastOneWord(result []models.FrontPost, words_wanted_list []string, t *testing.T) {
	for i := range result {
		content_list := strings.Split(result[i].Content, " ")
		found := false
		for _, word := range words_wanted_list {
			if slices.Contains(content_list, word) {
				found = true
				break
			}
		}
		assert.Equal(t, true, found, "At least one word should be in the content")
	}
}

func retweetAPost(post models.FrontPost, username, tokenRetweeterer string, r *gin.Engine, t *testing.T) models.FrontPost {
	retweetPost, _ := http.NewRequest("POST", "/twitsnap/retweet/"+post.Post_ID, nil)
	addAuthorization(retweetPost, tokenRetweeterer)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, retweetPost)

	retweet_post := models.FrontPost{}

	err := json.Unmarshal(first.Body.Bytes(), &retweet_post)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusCreated, first.Code, "Response should be 201")
	assert.Equal(t, true, retweet_post.User_Retweet, "User should have retweeted")
	assert.Equal(t, retweet_post.Content, post.Content, "Content should be the same")
	assert.Equal(t, retweet_post.Tags, post.Tags, "Tags should be the same")
	assert.Equal(t, retweet_post.Retweet_Author, username, "Retweet author should be the retweeter (Retweet)")
	assert.Equal(t, retweet_post.Retweets, post.Retweets, "Retweets should be the same")

	return retweet_post
}

func checkRetweetPost(post models.FrontPost, retweetAuthor string, t *testing.T) {
	assert.Equal(t, post.Retweet_Author, retweetAuthor, "Retweet author should be the retweeter (Retweet)")
	assert.Equal(t, post.Is_Retweet, true, "IsRetweet should be true")
}

func boookmarkPost(postID string, token string, r *gin.Engine, t *testing.T) {
	bookmarkPost, _ := http.NewRequest("POST", "/twitsnap/bookmark/"+postID, nil)
	addAuthorization(bookmarkPost, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, bookmarkPost)

	assert.Equal(t, http.StatusNoContent, first.Code)
}