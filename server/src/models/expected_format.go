package models

type PostExpectedFormat struct {
	Content string   `json:"content" validate:"required"`
	Public  bool     `json:"public"`
	Tags    []string `json:"tags" validate:"dive,required"`
	MediaURL string `json:"media_url"`
}

type LikeExpectedFormat struct {
	User_ID string `json:"user_id"`
}

type EditPostExpectedFormat struct {
	Content *string   `json:"content"`
	Tags    *[]string `json:"tags"`
	Public 	*bool     `json:"public"`
	MediaURL *string `json:"media_url"`
}

type ReturnPaginatedPosts struct {
	Data        []FrontPost `json:"data"`
	Pagination  Pagination  `json:"pagination"`
}

type FeedRequesst struct {
	FeedType     string `json:"feed_type"`
	WantedUserID string `json:"wanted_user_id"`
}
