package router

import (
	"server/src/controller"
	"server/src/database"

	"github.com/gin-gonic/gin"
)

func CreateRouter(db *database.Database) *gin.Engine {
	r := gin.Default()

	// create a postController
	postController := controller.NewPostController(db)

	r.POST("/twitsnap", func(c *gin.Context) {
		postController.NewPost(c)
	})

	r.POST("/twitsnap/like/:id", func(c *gin.Context) {
		// controller.PostLike(c)
	})

	r.PUT("/twitsnap/content/:id", func(c *gin.Context) {
		postController.UpdatePostContentByID(c, c.Param("id"))
	})

	r.PUT("/twitsnap/tags/:id", func(c *gin.Context) {
		postController.UpdatePostTagsByID(c, c.Param("id"))
	})

	r.GET("/twitsnap/feed", func(c *gin.Context) {
		postController.GetUserFeed(c)
	})

	r.GET("/twitsnap/myposts/:id", func(c *gin.Context) {
		// controller.GetAllTwitsnaps(c)
	})

	r.GET("/twitsnap/:id", func(c *gin.Context) {
		postController.GetPostByID(c, c.Param("id"))
	})

	r.GET("/twitsnap/interests", func(c *gin.Context) {
		postController.GetUserInterests(c)
	})

	r.DELETE("/twitsnap/:id", func(c *gin.Context) {
		postController.DeletePostByID(c, c.Param("id"))
	})

	return r
}
