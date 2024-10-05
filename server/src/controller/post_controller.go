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
	author_id, _ := context.Get("session_user_id")

	post, err := c.sv.FetchPostByID(postID, token.(string), author_id.(string))

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
	author_id, _ := context.Get("session_user_id")

	var editInfo models.EditPostExpectedFormat
	if err := context.ShouldBind(&editInfo); err != nil {
		_ = context.Error(postErrors.UnexpectedFormat())
		return
	}

	modPost, err := c.sv.ModifyPostByID(postID, editInfo, token.(string), author_id.(string))


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
		Pagination: models.Pagination{Limit: limitParams.Limit},
	}

	if hasMore {
		result.Pagination.Next_Offset =limitParams.Skip + limitParams.Limit
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) GetAllPosts(context *gin.Context) {
	token, _ := context.Get("tokenString")
	isUserAdmin, _ := context.Get("session_user_admin")

	if !isUserAdmin.(bool) {
		_ = context.Error(postErrors.AccssDenied())
		return
	}

	time := context.Query(TIME)
	skip := context.Query(SKIP)
	limit := context.Query(LIMIT)

	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.FetchAllPosts(limitParams, token.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := models.ReturnPaginatedPosts{
		Data: posts,
		Pagination: models.Pagination{Limit: limitParams.Limit},
	}

	if hasMore {
		result.Pagination.Next_Offset =limitParams.Skip + limitParams.Limit
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) HashtagsSearch(context *gin.Context) {
	token, _ := context.Get("tokenString")
	userID, _ := context.Get("session_user_id")


	hashtags := context.QueryArray(HASTAGS)
	time := context.Query(TIME)
	skip := context.Query(SKIP)
	limit := context.Query(LIMIT)

	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.FetchUserPostsByHashtags(hashtags, limitParams, userID.(string), token.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := models.ReturnPaginatedPosts{
		Data: posts,
		Pagination: models.Pagination{Limit: limitParams.Limit},
	}

	if hasMore {
		result.Pagination.Next_Offset =limitParams.Skip + limitParams.Limit
	}


	context.JSON(http.StatusOK, result)
}

func (c *PostController) WordsSearch(context *gin.Context) {
	token, _ := context.Get("tokenString")
	userID, _ := context.Get("session_user_id")

	words := context.Query(WORDS)

	time := context.Query(TIME)
	skip := context.Query(SKIP)
	limit := context.Query(LIMIT)

	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.WordsSearch(words, limitParams, userID.(string), token.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := models.ReturnPaginatedPosts{
		Data: posts,
		Pagination: models.Pagination{Limit: limitParams.Limit},
	}

	if hasMore {
		result.Pagination.Next_Offset =limitParams.Skip + limitParams.Limit
	}


	context.JSON(http.StatusOK, result)
}

func (c *PostController) LikePost(context *gin.Context) {
	postID := context.Param("id")
	userID, _ := context.Get("session_user_id")

	err := c.sv.LikePost(postID, userID.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func (c *PostController) UnLikePost(context *gin.Context) {
	postID := context.Param("id")
	userID, _ := context.Get("session_user_id")

	err := c.sv.UnLikePost(postID, userID.(string))

	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func (c *PostController) NoRoute(context *gin.Context) {
	_ = context.Error(postErrors.NotFound())
}