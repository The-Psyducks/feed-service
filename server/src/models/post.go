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
	Original_Author string `bson:"original_author"`
	Is_Retweet bool `bson:"is_retweet"`
}


func NewDBPost(author_id string, content string, tags []string, privacy bool) DBPost {
	return DBPost{
		Post_ID:   uuid.NewString(),
		Content:   content,
		Author_ID: author_id,
		Time:      time.Now().UTC(),
		Tags:      tags,
		Public:   privacy,
		Likes:   0,
		Original_Author: author_id,
		Is_Retweet: false,
	}
}

func NewRetweetDBPost(author_id string, content string, tags []string, privacy bool, original_author string) DBPost {
	return DBPost{
		Post_ID:   uuid.NewString(),
		Content:   content,
		Author_ID: author_id,
		Time:      time.Now().UTC(),
		Tags:      tags,
		Public:   privacy,
		Likes:   0,
		Original_Author: original_author,
		Is_Retweet: true,
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
	UserLiked  bool  `json:"liked_by_user"`
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