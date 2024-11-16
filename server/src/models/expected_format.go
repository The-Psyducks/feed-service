package models

type PostExpectedFormat struct {
	Content string   `json:"content" validate:"required"`
	Public  bool     `json:"public"`
	Mentions    []string `json:"mentions"`
	MediaInfo MediaInfo `json:"media_info"`
}

type LikeExpectedFormat struct {
	User_ID string `json:"user_id"`
}

type EditPostExpectedFormat struct {
	Content *string   `json:"content"`
	Public 	*bool     `json:"public"`
	MediaInfo *MediaInfo `json:"media_info"`
	Mentions    []string `json:"mentions"`
}

type ReturnPaginatedPosts struct {
	Data        []FrontPost `json:"data"`
	Pagination  Pagination  `json:"pagination"`
}

type FeedRequesst struct {
	FeedType     string `json:"feed_type"`
	WantedUserID string `json:"wanted_user_id"`
}
