package service

import (
	"server/src/models"
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

	return posts, hasMore, err
}