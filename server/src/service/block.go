package service

import (
	postErrors "server/src/all_errors"
)

func (c *Service) BlockPost(postID string) error {
	err := c.db.BlockPost(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}

func (c *Service) UnBlockPost(postID string) error {
	err := c.db.UnBlockPost(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}