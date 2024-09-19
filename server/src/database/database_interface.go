package database

import (
	"server/src/models"
)

type Database interface {
	AddNewPost(newPost models.DBPost) error
	GetPostByID(postID string) (models.DBPost, error)
	DeletePostByID(postID string) error
	EditPost(postID string, editInfo models.EditPostExpectedFormat) (models.DBPost, error)
	GetUserFeedFollowing(following []string) ([]models.DBPost, error)
	GetUserFeedInterests(interests []string) ([]models.DBPost, error)
	GetUserFeedSingle(userID string) ([]models.DBPost, error)
	WordSearchPosts(words string) ([]models.DBPost, error)
	GetUserHashtags(interests []string) ([]models.DBPost, error)
}
