package test

import (
	"bytes"
	"log"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	postErrors "server/src/all_errors"
	"server/src/auth"
	"server/src/models"
	"server/src/router"

	validator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {

	log.Println("TestNewPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "content #tag1 #tag2", Tags: []string{"tag1", "tag2"}, Mentions: []string{}, Public: true}
	req := newPostRequest(postBody)

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	addAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := models.FrontPost{}

	err = json.Unmarshal(first.Body.Bytes(), &result)

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusCreated, first.Code)
	makeResponseAsserions(t, http.StatusCreated, result, postBody, author_id, first.Code)
}

func TestNewPostWithMissInfo(t *testing.T) {

	log.Println("TestNewPostWithMissInfo")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "", Tags: []string{"tag1", "tag2"}, Public: true, Mentions: []string{}}
	req := newPostRequest(postBody)

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	addAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := postErrors.TwitSnapError{}

	err = json.Unmarshal(first.Body.Bytes(), &result)

	newPost := models.PostExpectedFormat{Content: "", Public: true}

	validate := validator.New()

	var errFMT postErrors.TwitSnapError
	if wrongFmtErr := validate.Struct(newPost); wrongFmtErr != nil {
		errFMT = postErrors.TwitSnapImportantFieldsMissing(wrongFmtErr)
	}

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, errFMT, result)
	assert.Equal(t, http.StatusBadRequest, first.Code)
}

func TestNewPostTooLong(t *testing.T) {

	log.Println("TestNewPostTooLong")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	postBody := PostBody{Content: "Psyduck, the perpetually confused Pokémon with its iconic yellow body and headache-induced psychic powers, waddles through life clutching its head, unintentionally unleashing bursts of immense energy, making it both endearingly clumsy and surprisingly powerful, a true enigma in the Pokémon world.", Tags: []string{}, Mentions: []string{}, Public: true}
	req := newPostRequest(postBody)

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	addAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := postErrors.TwitsnapTooLong()
	err = json.Unmarshal(first.Body.Bytes(), &result)

	errTooLong := postErrors.TwitsnapTooLong()

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusRequestEntityTooLarge, first.Code)
	assert.Equal(t, errTooLong.Detail, result.Detail)
}

func TestNewPostUnexpectedFormat(t *testing.T) {

	log.Println("TestNewPost")

	db := connectToDatabase()

	r := router.CreateRouter(db)

	author_id := "1234"
	post := struct {
		Content []string `json:"content"`
		T    []string `json:"tags"`
	}{
		Content: []string{"tag1", "tag2"},
		T:    []string{"tag1", "tag2"},
	}
	
	marshalledData, _ := json.Marshal(post)
	req, _ := http.NewRequest("POST", "/twitsnap", bytes.NewReader(marshalledData))

	req.Header.Add("content-type", "application/json")

	token, err := auth.GenerateToken(author_id, "username", true)

	if err != nil {
		log.Fatal("Error generating token: ", err)
	}

	addAuthorization(req, token)

	first := httptest.NewRecorder()
	r.ServeHTTP(first, req)

	result := postErrors.UnexpectedFormat()
	errUnexpectedFormat := postErrors.UnexpectedFormat()

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, http.StatusBadRequest, first.Code)
	assert.Equal(t, errUnexpectedFormat.Detail, result.Detail)
}