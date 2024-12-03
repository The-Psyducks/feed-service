package service

import (
	"log/slog"
	postErrors "server/src/all_errors"
	"server/src/models"
	"time"
)

func (c *Service) RetweetPost(postId string, userID string, token string) (*models.FrontPost, error) {
	post, err := c.db.GetPost(postId, userID)

	if err != nil {
		return nil, postErrors.TwitsnapNotFound(postId)
	}

	retweet := models.NewRetweetDBPost(post, userID)

	newRetweet, err := c.db.AddNewRetweet(retweet)

	if err != nil {
		return nil, postErrors.DatabaseError(err.Error())
	}

	newRetweet, err = addAuthorInfoToPost(newRetweet, token)

	if err != nil {
		return nil, postErrors.UserInfoError(err.Error())
	}

	slog.Info("Post retweeted: ", "post_id", postId, "User", userID, "time", time.Now())

	return &newRetweet, nil
}

func (c *Service) RemoveRetweet(postId string, userID string) error {
	err := c.db.DeleteRetweet(postId, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postId)
	}

	slog.Info("Retweet removed: ", "post_id", postId, "time", time.Now())

	return nil
}