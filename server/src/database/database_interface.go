package database

import (
	"server/src/models"
)

type Database interface {
	AddNewPost(newPost models.DBPost) (models.FrontPost, error)

	GetPost(postID string, askerID string) (models.FrontPost, error)

	DeletePost(postID string) error

	AddNewRetweet(newRetweet models.DBPost) (models.FrontPost, error)

	DeleteRetweet(postID string, userID string) error

	EditPost(postID string, editInfo models.EditPostExpectedFormat, askerID string) (models.FrontPost, error)

	GetAllPosts(limitConfig models.LimitConfig, askerID string) ([]models.FrontPost, bool, error)

	GetUserFeedFollowing(following []string, askerID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error)

	GetUserFeedInterests(interests []string, following []string, askerID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error)

	GetUserFeedSingle(userID string, limitConfig models.LimitConfig, askerID string, following []string) ([]models.FrontPost, bool, error)

	GetUserFeedRetweet(userID string, limitConfig models.LimitConfig, askerID string, following []string) ([]models.FrontPost, bool, error)

	WordSearchPosts(words string, following []string, askerID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error)

	GetUserHashtags(hashtags []string, following []string, askerID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error)

	GetUserMetrics(userID string, limits models.MetricLimits) (models.Metrics, error)

	LikeAPost(postID string, likerID string) error

	UnLikeAPost(postID string, likerID string) error

	AddFavorite(postID string, userID string) error
	
	RemoveFavorite(postID string, userID string) error

	GetUserFavorites(userID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error)

	BlockPost(postID string) error

	UnBlockPost(postID string) error

	ClearDB() error
}
