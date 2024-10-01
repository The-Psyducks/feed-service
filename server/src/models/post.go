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
	Likes   int  `bson:"likes"`
}


func NewDBPost(author_id string, content string, tags []string, privacy bool) DBPost {
	return DBPost{
		Post_ID:   uuid.NewString(),
		Content:   content,
		Author_ID: author_id,
		Time:      time.Now().UTC(),
		Tags:      tags,
		Public:   privacy,
	}
}

type AuthorInfo struct {
	Author_ID string `json:"author_id"`
	Username  string `json:"username"`
	Alias	 string `json:"alias"`
	PthotoURL    string `json:"photourl"`
}

type FrontPost struct {
	Post_ID   string    `json:"post_id"`
	Content   string    `json:"content"`
	Author_Info AuthorInfo    `json:"author"`
	Time      string `json:"time"`
	Public   bool    `json:"public"`
	Tags     []string  `json:"tags"`
	Likes   int  `json:"likes"`
	UserLiked  bool  `json:"user_liked"`
}

func NewFrontPost(post DBPost, author AuthorInfo, liked bool) FrontPost {
	return FrontPost{
		Post_ID:   post.Post_ID,
		Content:   post.Content,
		Author_Info: author,
		Time:      post.Time.Format(time.RFC3339),
		Tags:      post.Tags,
		Public:   post.Public,
		Likes:   post.Likes,
		UserLiked:  liked,
	}
}