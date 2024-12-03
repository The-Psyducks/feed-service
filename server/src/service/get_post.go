package service

import (
	"log/slog"
	postErrors "server/src/all_errors"
	"server/src/models"
	"time"
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

	slog.Info("Post retrieved: ", "post_id", postID, "User", userID, "time", time.Now())

	return &post, nil
}

func (c *Service) RemovePostByID(postID string) error {
	err := c.db.DeletePost(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	slog.Info("Post removed: ", "post_id", postID, "time", time.Now())

	return nil
}