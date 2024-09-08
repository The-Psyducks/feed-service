package post

import (
	"time"

	"github.com/google/uuid"
)

type DBPost struct {
	Post_ID   string    `bson:"post_id"`
	Content   string    `bson:"content"`
	Author_ID string    `bson:"author_id"`
	Time      time.Time `bson:"time"`
	// Privacy   int    `bson:"privacy"`
	Tags     []string  `bson:"tags"`
	// Likes is a list with the user_id of the users that liked the post
}

type PostOutput struct {
	Post_ID   string `json:"post_id"`
	Content   string `json:"content"`
	Author_ID string `json:"author_id"`
	Tags    []string `json:"tags"`
	Time      time.Time `json:"-"`
}

func NewDBPost(author_id string, content string, tags []string) DBPost {
	return DBPost{
		Post_ID:   uuid.NewString(),
		Content:   content,
		Author_ID: author_id,
		Time:      time.Now(),
		Tags:      tags,
	}
}

func NewPostOutput(author_id string, content string, post_id string, postTime time.Time) PostOutput {
	return PostOutput{
		Post_ID:   post_id,
		Content:   content,
		Author_ID: author_id,
		Time:      postTime,
	}
}
