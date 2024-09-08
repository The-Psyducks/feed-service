package post

type PostExpectedFromat struct {
	Content     string 	  `json:"content" validate:"required"`
	Author_ID	string    `json:"author_id" validate:"required"`
	// Privacy		int    `json:"privacy" validate:"required"`
	Tags		[]string `json:"tags" validate:"dive,required"`
}

type LikeExpectedFormat struct {
	User_ID		string    `json:"user_id"`
}

type EditPostExpectedFormat struct {
	Content		string    `json:"content" validate:"required"`
}

type EditPostTagsExpectedFormat struct {
	Tags		[]string    `json:"tags" validate:"required"`
}