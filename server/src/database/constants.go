package database

const (
	DATABASE_NAME       = "feed"
	FEED_COLLECTION     = "posts"
	LIKES_COLLECTION    = "likes"
	RETWEET_COLLECTION  = "retweets"
	BOOKMARK_COLLECTION = "bookmarks"
	TWTMETRICS_COLLECTION = "twtmetrics"
	TAGMETRICS_COLLECTION = "tagmetrics"
)

const (
	POST_ID_FIELD          = "post_id"
	CONTENT_FIELD          = "content"
	AUTHOR_ID_FIELD        = "author_id"
	TIME_FIELD             = "time"
	TAGS_FIELD             = "tags"
	LIKES_FIELD            = "likes"
	PUBLIC_FIELD           = "public"
	LIKERS_FIELD           = "likers"
	ORIGINAL_AUTHOR_FIELD  = "original_author"
	IS_RETWEET_FIELD       = "is_retweet"
	ORIGINAL_POST_ID_FIELD = "original_post_id"
	RETWEET_FIELD          = "retweets"
	RETWEETERS_FIELD       = "retweeters"
	RETWEET_AUTHOR_FIELD   = "retweet_author"
	MEDIA_INFO_FIELD        = "media_info"
	BOOKMARK_FIELD         = "bookmark"
	MENTIONS_FIELD         = "mentions"
	BLOCKED_FIELD		  = "blocked"
)

const (
	TOTAL_TWEETS = "total_tweets"
	HOURLY_FRECUENCY = "hourly_frecuency"
	TREND = "trend"
	LAST_UPDATED = "last_updated"
	DAY = "day"
)

const (
	ADMIN = "admin"
)