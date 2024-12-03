package service

import (
	"log/slog"
	"time"
	postErrors "server/src/all_errors"
)

func (c *Service) LikePost(postID string, userID string) error {
	err := c.db.LikeAPost(postID, userID)

	return err
}

func (c *Service) UnLikePost(postID string, userID string) error {
	err := c.db.UnLikeAPost(postID, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	slog.Info("Post unliked: ", "post_id", postID, "user", userID, "time", time.Now())

	return nil
}