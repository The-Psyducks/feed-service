package database

import (
	"server/src/post"
)

type Database interface {
	AddNewPost(newPost post.DBPost) error
	GetPostByID(postID string) (post.DBPost, error) 
	DeletePostByID(postID string) error
	UpdatePostContent(postID string, newContent string) (post.DBPost, error)
	UpdatePostTags(postID string, newTags []string) (post.DBPost, error) 
	GetUserFeed(following []string) ([]post.DBPost, error)
	WordSearchPosts(words string) ([]post.DBPost, error)
	GetUserInterests(interests []string) ([]post.DBPost, error)
}