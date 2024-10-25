package router

import (
	"server/src/database"
	"server/src/controller"
	"github.com/gin-gonic/gin"
	"server/src/middleware"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
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

	r.NoRoute(postController.NoRoute)

	return r
}

func setNewRelicConnection(r *gin.Engine) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("feed-micro"),
		newrelic.ConfigLicense("1133c491dc667d55d64f81f50324285dFFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	  )	  
	if err != nil {
		panic(fmt.Errorf("error connecting with new relic: %w",err))
	}
	r.Use(nrgin.Middleware(app))
}
