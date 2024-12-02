package service

import (
	postErrors "server/src/all_errors"
	"server/src/models"
)


func (c *Service) FetchUserFeed(feedRequest *models.FeedRequesst, user_id string, limitConfig models.LimitConfig, token string) ([]models.FrontPost, bool, error) {
	switch feedRequest.FeedType {
	case FOLLOWING:
		return c.fetchFollowingFeed(limitConfig, user_id, token)
	case FORYOU:
		return c.fetchForyouFeed(limitConfig, user_id, token)
	case SINGLE:
		return c.fetchForyouSingle(limitConfig, feedRequest.WantedUserID, user_id, token)
	case RETWEET:
		return c.fetchRetweetFeed(limitConfig, feedRequest.WantedUserID, user_id, token)
	}
	return []models.FrontPost{}, false, postErrors.BadFeedRequest(feedRequest.FeedType)
}

func (c *Service) fetchFollowingFeed(limitConfig models.LimitConfig, userID string, token string) ([]models.FrontPost, bool, error) {
	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	posts, hasMore, err := c.db.GetUserFeedFollowing(following, userID, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}

	posts, err = addAuthorInfoToPosts(posts, token)
	return posts, hasMore, err
}

func (c *Service) fetchForyouFeed(limitConfig models.LimitConfig, userID string, token string) ([]models.FrontPost, bool, error) {

	interests, err := getUserInterestsWp(userID, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}

	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	posts, hasMore, err := c.db.GetUserFeedInterests(interests, following, userID, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}

	posts, err = addAuthorInfoToPosts(posts, token)
	return posts, hasMore, err
}

func (c *Service) fetchForyouSingle(limitConfig models.LimitConfig, wantedUserID string, userID string, token string) ([]models.FrontPost, bool, error) {

	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if wantedUserID == userID {
		following = append(following, userID)
	}

	posts, hasMore, err := c.db.GetUserFeedSingle(wantedUserID, limitConfig, userID, following)
	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}

	posts, err = addAuthorInfoToPosts(posts, token)
	return posts, hasMore, err
}

func (c *Service) fetchRetweetFeed(limitConfig models.LimitConfig, wantedUserID string, userID string, token string) ([]models.FrontPost, bool, error) {

	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}

	posts, hasMore, err := c.db.GetUserFeedRetweet(wantedUserID, limitConfig, userID, following)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}

	posts, err = addAuthorInfoToPosts(posts, token)
	return posts, hasMore, err
}