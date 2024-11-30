package router

import (
	"fmt"
	"log"
	"server/src/controller"
	"server/src/database"
	"server/src/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func CreateRouter(db database.Database) *gin.Engine {
	r := gin.Default()

	addCorsConfiguration(r)
	setNewRelicConnection(r)

	r.Use(middleware.ErrorManager())
	r.Use(middleware.AuthMiddleware())

	amqp, err := CreateProducer()
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}

	postController := controller.NewPostController(db, amqp)

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

	r.GET("/twitsnap/metrics", postController.GetUserMetrics)

	r.GET("/twitsnap/trending", postController.GetTrendingTopics)

	r.POST("/twitsnap/block/:id", postController.BlockPost)

	r.DELETE("/twitsnap/block/:id", postController.UnBlockPost)

	r.NoRoute(postController.NoRoute)

	return r
}

func addCorsConfiguration(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	r.Use(cors.New(config))
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

