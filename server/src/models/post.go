package models

import (
	"time"

	"github.com/google/uuid"
)

type DBPost struct {
	Post_ID         string    `bson:"post_id"`
	Content         string    `bson:"content"`
	Author_ID       string    `bson:"author_id"`
	Time            time.Time `bson:"time"`
	Public          bool      `bson:"public"`
	Tags            []string  `bson:"tags"`
	Likes           int       `bson:"likes"`
	Retweets        int       `bson:"retweets"`
	IsRetweet       bool      `bson:"is_retweet"`
	OriginalPostID  string    `bson:"original_post_id"`
	RetweetAuthorID string    `bson:"retweet_author"`
	MediaURL        string    `bson:"media_url"`
}

func NewDBPost(author_id string, content string, tags []string, privacy bool, mediaUrl string) DBPost {
	postID := uuid.NewString()
	return DBPost{
		Post_ID:         postID,
		Content:         content,
		Author_ID:       author_id,
		Time:            time.Now().UTC(),
		Tags:            tags,
		Public:          privacy,
		Likes:           0,
		Retweets:        0,
		OriginalPostID:  postID,
		RetweetAuthorID: author_id,
		IsRetweet:       false,
		MediaURL:        mediaUrl,
	}
}

func NewRetweetDBPost(post FrontPost, author_id string) DBPost {
	return DBPost{
		Post_ID:         uuid.NewString(),
		Content:         post.Content,
		Author_ID:       post.Author_Info.Author_ID,
		Time:            time.Now().UTC(),
		Tags:            post.Tags,
		Public:          post.Public,
		Likes:           0,
		Retweets:        0,
		RetweetAuthorID: author_id,
		OriginalPostID:  post.OriginalPostID,
		IsRetweet:       true,
		MediaURL:        post.MediaURL,
	}
}

type AuthorInfo struct {
	Author_ID string `json:"author_id"`
	Username  string `json:"username"`
	Alias     string `json:"alias"`
	PthotoURL string `json:"photourl"`
}

type FrontPost struct {
	Post_ID     string     `json:"post_id"`
	Content     string     `json:"content"`
	Author_Info AuthorInfo `json:"author"`
	Time        string     `json:"time"`
	Public      bool       `json:"public"`
	Tags        []string   `json:"tags"`
	Likes       int        `json:"likes"`
	Retweets    int        `json:"retweets"`
	UserLiked   bool       `json:"liked_by_user"`

	IsRetweet      bool   `bson:"is_retweet"`
	OriginalPostID string `bson:"original_post_id"`
	RetweetAuthor  string `bson:"retweet_author"`
	MediaURL       string `bson:"media_url"`
}

func NewFrontPost(post DBPost, author AuthorInfo, liked bool) FrontPost {
	return FrontPost{
		Post_ID:        post.Post_ID,
		Content:        post.Content,
		Author_Info:    author,
		Time:           post.Time.Format(time.RFC3339),
		Tags:           post.Tags,
		Public:         post.Public,
		Likes:          post.Likes,
		Retweets:       post.Retweets,
		UserLiked:      liked,
		MediaURL:       post.MediaURL,
		OriginalPostID: post.OriginalPostID,
		IsRetweet:      post.IsRetweet,
		RetweetAuthor:  post.RetweetAuthorID,
	}

}
