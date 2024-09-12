package router

import (
	"server/src/controller"
	"server/src/database"

	"github.com/gin-gonic/gin"
)

func CreateRouter(db database.Database) *gin.Engine {
	r := gin.Default()

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

	r.GET("/twitsnap/interests", func(c *gin.Context) {
		postController.GetUserInterests(c)
	})

	r.GET("/twitsnap/wordsearch", func(c *gin.Context) {
		postController.WordsSearch(c)
	})

	r.DELETE("/twitsnap/:id", func(c *gin.Context) {
		postController.DeletePostByID(c, c.Param("id"))
	})

	return r
}
