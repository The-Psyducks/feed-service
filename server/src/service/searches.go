package service

import (
	"log/slog"
	"server/src/models"
	"time"
)

func (c *Service) FetchUserPostsByHashtags(hashtags []string, limitConfig models.LimitConfig, userID string, token string) ([]models.FrontPost, bool, error) {

	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}

	posts, hasMore, err := c.db.GetUserHashtags(hashtags, following, userID, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}

	posts, err = addAuthorInfoToPosts(posts, token)

	slog.Info("Hashtags feed retrieved: ", "user_id", userID, "time", time.Now(), "count", len(posts), "hashtags", hashtags)

	return posts, hasMore, err
}

func (c *Service) WordsSearch(words string, limitConfig models.LimitConfig, userID string, token string) ([]models.FrontPost, bool, error) {
	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	posts, hasMore, err := c.db.WordSearchPosts(words, following, userID, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}

	posts, err = addAuthorInfoToPosts(posts, token)

	slog.Info("Words search feed retrieved: ", "user_id", userID, "time", time.Now(), "count", len(posts), "words", words)

	return posts, hasMore, err
}