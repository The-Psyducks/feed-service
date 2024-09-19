package router

import (
	"server/src/database"
	"server/src/controller"
	"github.com/gin-gonic/gin"
	"server/src/middleware"
)

func CreateRouter(db database.Database) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ErrorManager())

	postController := controller.NewPostController(db)

	r.POST("/twitsnap", postController.NewPost)

	r.PUT("/twitsnap/edit/:id", postController.UpdatePostByID)

	r.GET("/twitsnap/feed/:id", postController.GetUserFeed)

	r.GET("/twitsnap/:id", postController.GetPostByID)

	r.GET("/twitsnap/hashtags", postController.HashtagsSearch)

	r.GET("/twitsnap/wordsearch",postController.WordsSearch)

	r.DELETE("/twitsnap/:id", postController.DeletePostByID)

	return r
}
