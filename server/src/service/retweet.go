package service

import (
	postErrors "server/src/all_errors"
	"server/src/database"
	"server/src/models"
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

	return &newRetweet, nil
}

func (c *Service) RemoveRetweet(postId string, userID string) error {
	err := c.db.DeleteRetweet(postId, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postId)
	}

	return nil
}

func (c *Service) FetchAllPosts(limitConfig models.LimitConfig, token string) ([]models.FrontPost, bool, error) {

	posts, hasMore, err := c.db.GetAllPosts(limitConfig, database.ADMIN)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}

	posts, err = addAuthorInfoToPosts(posts, token)

	return posts, hasMore, err
}