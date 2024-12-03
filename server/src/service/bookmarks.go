package service

import (
	"log/slog"
	postErrors "server/src/all_errors"
	"server/src/models"
	"time"
)

func (c *Service) BookmarkPost(postID string, userID string) error {
	err := c.db.AddFavorite(postID, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	slog.Info("Post bookmarked: ", "post_id", postID, "time", time.Now())

	return nil
}

func (c *Service) UnBookmarkPost(postID string, userID string) error {
	err := c.db.RemoveFavorite(postID, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	slog.Info("Post unbookmarked: ", "post_id", postID, "time", time.Now())

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

	slog.Info("User favorites retrieved: ", "user_id", userID, "time", time.Now(), "count", len(posts))

	return posts, hasMore, nil
}
