package database

import (
	"server/src/post"
)

type Database interface {
	AddNewPost(newPost post.DBPost) error
	GetPostByID(postID string) (post.DBPost, error) 
	DeletePostByID(postID string) error
	EditPost(postID string, editInfo post.EditPostExpectedFormat) (post.DBPost, error)
	GetUserFeed(following []string) ([]post.DBPost, error)
	WordSearchPosts(words string) ([]post.DBPost, error)
	GetUserInterests(interests []string) ([]post.DBPost, error)
}