package database

import (
	"server/src/models"
)

type Database interface {
	AddNewPost(newPost models.DBPost) (models.FrontPost, error)

	GetPostByID(postID string) (models.FrontPost, error)

	DeletePostByID(postID string) error

	EditPost(postID string, editInfo models.EditPostExpectedFormat) (models.FrontPost, error)

	GetUserFeedFollowing(following []string,  limitConfig models.LimitConfig) ([]models.FrontPost, error)

	GetUserFeedInterests(interests []string, following []string, limitConfig models.LimitConfig) ([]models.FrontPost, error)

	GetUserFeedSingle(userID string, limitConfig models.LimitConfig) ([]models.FrontPost, error)

	WordSearchPosts(words string, following []string, limitConfig models.LimitConfig) ([]models.FrontPost, error)

	GetUserHashtags(hashtags []string, following []string, limitConfig models.LimitConfig) ([]models.FrontPost, error)

	LikeAPost(postID string) error

	UnLikeAPost(postID string) error

}
