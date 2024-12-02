package service

import (
	postErrors "server/src/all_errors"
	"server/src/models"
)

func (c *Service) BookmarkPost(postID string, userID string) error {
	err := c.db.AddFavorite(postID, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}

func (c *Service) UnBookmarkPost(postID string, userID string) error {
	err := c.db.RemoveFavorite(postID, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}

func (c *Service) GetUserFavorites(userID string, limitiConfig models.LimitConfig, token string) ([]models.FrontPost, bool, error) {
	bookmarks, hasMore, err := c.db.GetUserFavorites(userID, limitiConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	posts, err := addAuthorInfoToPosts(bookmarks, token)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	return posts, hasMore, nil
}
