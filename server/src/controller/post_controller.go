package controller

import (
	"net/http"
	postErrors "server/src/all_errors"
	"server/src/database"
	"server/src/models"
	"server/src/service"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	sv *service.Service
}

func NewPostController(sv database.Database) *PostController {
	return &PostController{sv: service.NewService(sv)}
}

func (c *PostController) NewPost(context *gin.Context) {

	var newPost models.PostExpectedFormat
	if err := context.ShouldBind(&newPost); err != nil {
		_ = context.Error(postErrors.UnexpectedFormat())
		return
	}

	postNew, err := c.sv.CreatePost(&newPost)

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := gin.H{
		"post": postNew,
	}

	context.JSON(http.StatusCreated, result)
}

func (c *PostController) GetPostByID(context *gin.Context) {

	postID := context.Param("id")

	post, err := c.sv.FetchPostByID(postID)

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := gin.H{
		"post": post,
	}
	context.JSON(http.StatusOK, result)
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

	var editInfo models.EditPostExpectedFormat
	if err := context.ShouldBind(&editInfo); err != nil {
		_ = context.Error(postErrors.UnexpectedFormat())
		return
	}

	modPost, err := c.sv.ModifyPostByID(postID, editInfo)

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := gin.H{
		"post": modPost,
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) GetUserFeed(context *gin.Context) {

	username := context.Param("username")

	feedType := context.Query("feed_type")
	time := context.Query("from_time")
	skip := context.Query("skip")
	limit := context.Query("limit")

	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.FetchUserFeed(username, feedType, limitParams)

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
	hashtags := context.QueryArray("tags")

	time := context.Query("from_time")
	skip := context.Query("skip")
	limit := context.Query("limit")

	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.FetchUserPostsByHashtags(hashtags, limitParams)

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
	words := context.Query("words")

	time := context.Query("from_time")
	skip := context.Query("skip")
	limit := context.Query("limit")

	limitParams := models.NewLimitConfig(time, skip, limit)

	posts, hasMore, err := c.sv.WordsSearch(words, limitParams)

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
