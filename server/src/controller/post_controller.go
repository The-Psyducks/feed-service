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
		context.Error(postErrors.UnexpectedFormat())
		return
	}

	postNew, err := c.sv.NewPost(&newPost)

	if err != nil {
		context.Error(err)
		return
	}

	result := gin.H{
		"post": postNew,
	}

	context.JSON(http.StatusCreated, result)
}

func (c *PostController) GetPostByID(context *gin.Context, postID string) {

	post, err := c.sv.GetPostByID(postID)

	if err != nil {
		context.Error(err)
		return
	}

	result := gin.H{
		"post": post,
	}
	context.JSON(http.StatusOK, result)
}

func (c *PostController) DeletePostByID(context *gin.Context, postID string) {
	err := c.sv.DeletePostByID(postID)

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func (c *PostController) UpdatePostByID(context *gin.Context, postID string) {
	var editInfo models.EditPostExpectedFormat
	if err := context.ShouldBind(&editInfo); err != nil {
		context.JSON(http.StatusBadRequest, postErrors.UnexpectedFormat())
		return
	}

	modPost, err := c.sv.UpdatePostByID(postID, editInfo)

	if err != nil {
		context.Error(err)
		return
	}

	result := gin.H{
		"post": modPost,
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) GetUserFeed(context *gin.Context) {
	following := context.QueryArray("following")

	posts, err := c.sv.GetUserFeed(following)

	if err != nil {
		context.Error(err)
		return
	}

	result := gin.H{
		"posts": posts,
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) GetUserPostsByHashtags(context *gin.Context) {
	hashtags := context.QueryArray("tags")

	posts, err := c.sv.GetUserPostsByHashtags(hashtags)

	if err != nil {
		context.Error(err)
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
		context.Error(err)
		return
	}

	result := gin.H{
		"posts": posts,
	}

	context.JSON(http.StatusOK, result)
}
