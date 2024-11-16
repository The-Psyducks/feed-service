package models

type MentionNotificationRequest struct {
	UserId   string `json:"user_id"`
	TaggerId string `json:"tagger_id"`
	PostId   string `json:"post_id"`
}