package service

import (
	"log/slog"
	postErrors "server/src/all_errors"
	"time"
)

func (c *Service) BlockPost(postID string) error {
	err := c.db.BlockPost(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	slog.Info("Post blocked: ", "post_id", postID, "time", time.Now())

	return nil
}

func (c *Service) UnBlockPost(postID string) error {
	err := c.db.UnBlockPost(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	slog.Info("Post unblocked: ", "post_id", postID, "time", time.Now())

	return nil
}