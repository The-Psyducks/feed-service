package models

import (
	"time"

	"github.com/google/uuid"
)

type DBPost struct {
	Post_ID   string    `bson:"post_id"`
	Content   string    `bson:"content"`
	Author_ID string    `bson:"author_id"`
	Time      time.Time `bson:"time"`
	Public   bool    `bson:"public"`
	Tags     []string  `bson:"tags"`
	// Likes is a list with the user_id of the users that liked the post
}


func NewDBPost(author_id string, content string, tags []string, privacy bool) DBPost {
	return DBPost{
		Post_ID:   uuid.NewString(),
		Content:   content,
		Author_ID: author_id,
		Time:      time.Now(),
		Tags:      tags,
		Public:   privacy,
	}
}
