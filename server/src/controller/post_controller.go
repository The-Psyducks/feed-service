package controller

import (
	"errors"
	"net/http"
	postErrors "server/src/all_errors"
	"server/src/database"
	"server/src/models"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

type PostController struct {
	db database.Database
}

func NewPostController(db database.Database) *PostController {
	return &PostController{db: db}
}

func (c *PostController) NewPost(context *gin.Context) {
	var newPost models.PostExpectedFormat
	if err := context.ShouldBind(&newPost); err != nil {
		context.JSON(http.StatusBadRequest, postErrors.UnexpectedFormat())
		return
	}

	validate := validator.New()
	if err := validate.Struct(newPost); err != nil {
		context.JSON(http.StatusBadRequest, postErrors.TwitSnapImportantFieldsMissing(err))
		return
	}

	if len(newPost.Content) > 280 {
		context.JSON(http.StatusRequestEntityTooLarge, postErrors.TwitsnapTooLong())
		return
	}

	postNew := models.NewDBPost(newPost.Author_ID, newPost.Content, newPost.Tags, newPost.Public)

	if err := c.db.AddNewPost(postNew); err != nil {
		context.JSON(http.StatusInternalServerError, postErrors.DatabaseError())
		return
	}

	result := gin.H{
		"post": postNew,
	}

	context.JSON(http.StatusCreated, result)
}

func (c *PostController) GetPostByID(context *gin.Context, postID string) {

	post, err := c.db.GetPostByID(postID)

	if err != nil {
		context.JSON(http.StatusNotFound, postErrors.TwitsnapNotFound(postID))
		return
	}

	result := gin.H{
		"post": post,
	}
	context.JSON(http.StatusOK, result)
}

func (c *PostController) DeletePostByID(context *gin.Context, postID string) {
	err := c.db.DeletePostByID(postID)

	if err != nil {
		context.JSON(http.StatusNotFound, postErrors.TwitsnapNotFound(postID))
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

	validate := validator.New()
	if err := validate.Struct(editInfo); err != nil {
		context.JSON(http.StatusBadRequest, postErrors.TwitSnapImportantFieldsMissing(err))
		return
	}

	modPost, err := c.db.EditPost(postID, editInfo)

	if err != nil {
		if errors.Is(err, postErrors.ErrTwitsnapNotFound) {
			context.JSON(http.StatusNotFound, postErrors.TwitsnapNotFound(postID))
			return
		} else {
			context.JSON(http.StatusInternalServerError, postErrors.DatabaseError())
			return
		}
	}

	result := gin.H{
		"post": modPost,
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) GetUserFeed(context *gin.Context) {
	following := context.QueryArray("following")

	posts, err := c.db.GetUserFeed(following)

	if err != nil {
		context.JSON(http.StatusNotFound, err.Error())
		return
	}

	result := gin.H{
		"posts": posts,
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) GetUserPostsByHashtags(context *gin.Context) {
	interests := context.QueryArray("tags")

	posts, err := c.db.GetUserHashtags(interests)

	if err != nil {
		context.JSON(http.StatusNotFound, err.Error())
		return
	}

	if posts == nil {
		context.JSON(http.StatusNotFound, postErrors.NoTagsFound())
		return
	}

	result := gin.H{
		"posts": posts,
	}

	context.JSON(http.StatusOK, result)
}

func (c *PostController) WordsSearch(context *gin.Context) {
	words := context.Query("words")

	posts, err := c.db.WordSearchPosts(words)

	if err != nil {
		context.JSON(http.StatusNotFound, err.Error())
		return
	}

	if posts == nil {
		context.JSON(http.StatusNotFound, postErrors.NoWordssFound())
		return
	}

	result := gin.H{
		"posts": posts,
	}

	context.JSON(http.StatusOK, result)
}
