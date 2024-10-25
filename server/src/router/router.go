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
	r.Use(middleware.AuthMiddleware())

	postController := controller.NewPostController(db)

	r.POST("/twitsnap", postController.NewPost)
	
	r.POST("/twitsnap/retweet/:id", postController.NewPostRetweet)

	r.DELETE("/twitsnap/retweet/:id", postController.DeleteRetweet)
	
	r.PUT("/twitsnap/edit/:id", postController.UpdatePostByID)

	r.GET("/twitsnap/feed", postController.GetUserFeed)

	r.GET("/twitsnap/:id", postController.GetPostByID)

	r.GET("/twitsnap/hashtag-search", postController.HashtagsSearch)

	r.GET("/twitsnap/word-search",postController.WordsSearch)

	r.DELETE("/twitsnap/:id", postController.DeletePostByID)

	r.POST("/twitsnap/like/:id", postController.LikePost)

	r.DELETE("/twitsnap/like/:id", postController.UnLikePost)

	r.GET("/twitsnap/all", postController.GetAllPosts)

	r.POST("/twitsnap/bookmark/:id", postController.BookmarkPost)

	r.DELETE("/twitsnap/bookmark/:id", postController.UnBookmarkPost)

	r.GET("/twitsnap/bookmarks", postController.GetBookmarks)

	r.NoRoute(postController.NoRoute)

	return r
}
