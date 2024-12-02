package models

const (
	NEW_CONTENT = "NEW_CONTENT"
)

type UserMetrics struct {
	Likes    int `json:"likes"`
	Retweets int `json:"retweets"`
	Posts    int `json:"posts"`
}

type QueueMessage struct {
	MessageType string `json:"message_type"`
	Message     interface{} `json:"message"`
}

type New_Content struct {
	Hashtags []string `json:"hashtags"`
	Timestamp string `json:"timestamp"`
}