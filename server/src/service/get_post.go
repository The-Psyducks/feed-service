package service

import (
	postErrors "server/src/all_errors"
	"server/src/models"
)

func (c *Service) FetchPostByID(postID string, token string, userID string) (*models.FrontPost, error) {

	post, err := c.db.GetPost(postID, userID)

	if err != nil {
		return nil, postErrors.TwitsnapNotFound(postID)
	}

	post, err = addAuthorInfoToPost(post, token)

	if err != nil {
		return nil, postErrors.UserInfoError(err.Error())
	}

	return &post, nil
}

func (c *Service) RemovePostByID(postID string) error {
	err := c.db.DeletePost(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}