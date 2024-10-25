package models

import (
	"time"

	"github.com/google/uuid"
)

type DBPost struct {
	Post_ID           string    `bson:"post_id"`
	Content           string    `bson:"content"`
	Author_ID         string    `bson:"author_id"`
	Time              time.Time `bson:"time"`
	Public            bool      `bson:"public"`
	Tags              []string  `bson:"tags"`
	Likes             int       `bson:"likes"`
	Retweets          int       `bson:"retweets"`
	Is_Retweet        bool      `bson:"is_retweet"`
	Original_Post_ID  string    `bson:"original_post_id"`
	Retweet_Author_ID string    `bson:"retweet_author"`
	Media_URL          string    `bson:"media_url"`
}

func NewDBPost(author_id string, content string, tags []string, privacy bool, mediaUrl string) DBPost {
	postID := uuid.NewString()
	return DBPost{
		Post_ID:           postID,
		Content:           content,
		Author_ID:         author_id,
		Time:              time.Now().UTC(),
		Tags:              tags,
		Public:            privacy,
		Likes:             0,
		Retweets:          0,
		Original_Post_ID:  postID,
		Retweet_Author_ID: author_id,
		Is_Retweet:        false,
		Media_URL:          mediaUrl,
	}
}

func NewRetweetDBPost(post FrontPost, author_id string) DBPost {
	return DBPost{
		Post_ID:           uuid.NewString(),
		Content:           post.Content,
		Author_ID:         post.Author_Info.Author_ID,
		Time:              time.Now().UTC(),
		Tags:              post.Tags,
		Public:            post.Public,
		Likes:             post.Likes,
		Retweets:          post.Retweets,
		Retweet_Author_ID: author_id,
		Original_Post_ID:  post.Original_Post_ID,
		Is_Retweet:        true,
		Media_URL:          post.Media_URL,
	}
}

type AuthorInfo struct {
	Author_ID string `json:"author_id"`
	Username  string `json:"username"`
	Alias     string `json:"alias"`
	PthotoURL string `json:"photo_url"`
}

type FrontPost struct {
	Post_ID          string     `json:"post_id"`
	Content          string     `json:"content"`
	Author_Info      AuthorInfo `json:"author"`
	Time             string     `json:"time"`
	Public           bool       `json:"public"`
	Tags             []string   `json:"tags"`
	Likes            int        `json:"likes"`
	Retweets         int        `json:"retweets"`
	User_Liked       bool       `json:"user_liked"`
	User_Retweet     bool       `json:"user_retweet"`
	Is_Retweet       bool       `json:"is_retweet"`
	Original_Post_ID string     `json:"original_post_id"`
	Retweet_Author   string     `json:"retweet_author"`
	Media_URL        string     `json:"media_url"`
	Bookmark		 bool       `json:"bookmark"`
}

func NewFrontPost(post DBPost, author AuthorInfo, liked bool, retweeted bool, bookmarked bool) FrontPost {
	return FrontPost{
		Post_ID:          post.Post_ID,
		Content:          post.Content,
		Author_Info:      author,
		Time:             post.Time.Format(time.RFC3339Nano),
		Tags:             post.Tags,
		Public:           post.Public,
		Likes:            post.Likes,
		Retweets:         post.Retweets,
		User_Liked:       liked,
		User_Retweet:     retweeted,
		Media_URL:        post.Media_URL,
		Original_Post_ID: post.Original_Post_ID,
		Is_Retweet:       post.Is_Retweet,
		Retweet_Author:   post.Retweet_Author_ID,
		Bookmark:		  bookmarked,
	}

}
