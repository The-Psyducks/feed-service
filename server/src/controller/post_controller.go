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

	userID := context.Param("id")

	feedType := context.Query("feed_type")

	posts, err := c.sv.FetchUserFeed(userID, feedType)

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := gin.H{
		"posts": posts,
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) GetUserPostsByHashtags(context *gin.Context) {
	hashtags := context.QueryArray("tags")

	posts, err := c.sv.FetchUserPostsByHashtags(hashtags)

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := gin.H{
		"posts": posts,
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) WordsSearch(context *gin.Context) {
	words := context.Query("words")

	posts, err := c.sv.WordsSearch(words)

	if err != nil {
		_ = context.Error(err)
		return
	}

	result := gin.H{
		"posts": posts,
	}

	context.JSON(http.StatusOK, result)
}
