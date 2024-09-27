package controller

import (
	"net/http"
	postErrors "server/src/all_errors"
	"server/src/database"
	"server/src/models"
	"server/src/service"

	"github.com/gin-gonic/gin"
)

const (
	TIME = "time"
	SKIP = "skip"
	LIMIT = "limit"
	FEED = "feed_type"
	WORDS = "words"
	HASTAGS = "tags"
	WANTED_ID = "wanted_user_id"
)

type PostController struct {
	sv *service.Service
}

func NewPostController(sv database.Database) *PostController {
	return &PostController{sv: service.NewService(sv)}
}

func (c *PostController) NewPost(context *gin.Context) {

	token, _ := context.Get("tokenString")
	author_id, _ := context.Get("session_user_id")

	var newPost models.PostExpectedFormat
	if err := context.ShouldBind(&newPost); err != nil {
		_ = context.Error(postErrors.UnexpectedFormat())
		return
	}

	postNew, err := c.sv.CreatePost(&newPost, author_id.(string), token.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusCreated, postNew)
}

func (c *PostController) GetPostByID(context *gin.Context) {

	postID := context.Param("id")
	token, _ := context.Get("tokenString")

	post, err := c.sv.FetchPostByID(postID, token.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusOK, post)
}

func (c *PostController) DeletePostByID(context *gin.Context) {

	postID := context.Param("id")

	err := c.sv.RemovePostByID(postID)

	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func (c *PostController) UpdatePostByID(context *gin.Context) {
	
	postID := context.Param("id")
	token, _ := context.Get("tokenString")

	var editInfo models.EditPostExpectedFormat
	if err := context.ShouldBind(&editInfo); err != nil {
		_ = context.Error(postErrors.UnexpectedFormat())
		return
	}

	modPost, err := c.sv.ModifyPostByID(postID, editInfo, token.(string))


	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusOK, modPost)
}

func (c *PostController) GetUserFeed(context *gin.Context) {
	token, _ := context.Get("tokenString")
	author_id, _ := context.Get("session_user_id")
	
	time := context.Query(TIME)
	skip := context.Query(SKIP)
	limit := context.Query(LIMIT)
	feed_type := context.Query(FEED)
	wanted_id := context.Query(WANTED_ID)

	feedRequest := models.FeedRequesst{FeedType: feed_type, WantedUserID: wanted_id}


	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.FetchUserFeed(&feedRequest, author_id.(string), limitParams, token.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := models.ReturnPaginatedPosts{
		Data: posts,
		Limit: limitParams.Limit,
	}

	if hasMore {
		result.Next_Offset =limitParams.Skip + limitParams.Limit
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) HashtagsSearch(context *gin.Context) {
	username := context.Param("username")
	token, _ := context.Get("tokenString")

	hashtags := context.QueryArray(HASTAGS)
	time := context.Query(TIME)
	skip := context.Query(SKIP)
	limit := context.Query(LIMIT)

	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.FetchUserPostsByHashtags(hashtags, limitParams, username, token.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := models.ReturnPaginatedPosts{
		Data: posts,
		Limit: limitParams.Limit,
	}

	if hasMore {
		result.Next_Offset =limitParams.Skip + limitParams.Limit
	}


	context.JSON(http.StatusOK, result)
}

func (c *PostController) WordsSearch(context *gin.Context) {
	username := context.Param("username")
	token, _ := context.Get("tokenString")

	words := context.Query(WORDS)

	time := context.Query(TIME)
	skip := context.Query(SKIP)
	limit := context.Query(LIMIT)

	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.WordsSearch(words, limitParams, username, token.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := models.ReturnPaginatedPosts{
		Data: posts,
		Limit: limitParams.Limit,
	}

	if hasMore {
		result.Next_Offset =limitParams.Skip + limitParams.Limit
	}


	context.JSON(http.StatusOK, result)
}

func (c *PostController) LikePost(context *gin.Context) {
	postID := context.Param("id")

	err := c.sv.LikePost(postID)

	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func (c *PostController) UnLikePost(context *gin.Context) {
	postID := context.Param("id")

	err := c.sv.UnLikePost(postID)

	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}
