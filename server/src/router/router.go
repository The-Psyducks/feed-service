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

	// create a postController
	postController := controller.NewPostController(db)

	r.POST("/twitsnap", func(c *gin.Context) {
		postController.NewPost(c)
	})

	r.PUT("/twitsnap/edit/:id", func(c *gin.Context) {
		postController.UpdatePostByID(c, c.Param("id"))
	})

	r.GET("/twitsnap/feed", func(c *gin.Context) {
		postController.GetUserFeed(c)
	})

	r.GET("/twitsnap/:id", func(c *gin.Context) {
		postController.GetPostByID(c, c.Param("id"))
	})

	r.GET("/twitsnap/hashtags", func(c *gin.Context) {
		postController.GetUserPostsByHashtags(c)
	})

	r.GET("/twitsnap/wordsearch", func(c *gin.Context) {
		postController.WordsSearch(c)
	})

	r.DELETE("/twitsnap/:id", func(c *gin.Context) {
		postController.DeletePostByID(c, c.Param("id"))
	})

	return r
}
