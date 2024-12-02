package service

import (
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

	return nil
}