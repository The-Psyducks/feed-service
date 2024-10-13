package service

import (
	"errors"
	"log"
	postErrors "server/src/all_errors"
	"server/src/database"
	"server/src/models"

	validator "github.com/go-playground/validator/v10"
)

const (
	FOLLOWING = "following"
	FORYOU    = "foryou"
	SINGLE    = "single"
)

type Service struct {
	db database.Database
}

func NewService(db database.Database) *Service {
	return &Service{db: db}
}

func (c *Service) CreatePost(newPost *models.PostExpectedFormat, author_id string, token string) (*models.FrontPost, error) {

	validate := validator.New()
	if err := validate.Struct(newPost); err != nil {
		return nil, postErrors.TwitSnapImportantFieldsMissing(err)
	}

	if len(newPost.Content) > 280 {
		return nil, postErrors.TwitsnapTooLong()
	}

	postNew := models.NewDBPost(author_id, newPost.Content, newPost.Tags, newPost.Public, newPost.MediaURL)

	newPosted, err := c.db.AddNewPost(postNew)

	if err != nil {
		return nil, postErrors.DatabaseError(err.Error())
	}

	newPosted, err = addAuthorInfoToPost(newPosted, token)

	if err != nil {
		return nil, postErrors.UserInfoError(err.Error())
	}

	return &newPosted, nil
}

func (c *Service) FetchPostByID(postID string, token string, userID string) (*models.FrontPost, error) {

	post, err := c.db.GetPostByID(postID, userID)

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
	err := c.db.DeletePostByID(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}

func (c *Service) ModifyPostByID(postID string, editInfo models.EditPostExpectedFormat, token string, userID string) (*models.FrontPost, error) {
	validate := validator.New()
	if err := validate.Struct(editInfo); err != nil {
		return nil, postErrors.TwitSnapImportantFieldsMissing(err)
	}

	modPost, err := c.db.EditPost(postID, editInfo, userID)

	if err != nil {
		if errors.Is(err, postErrors.ErrTwitsnapNotFound) {
			return nil, postErrors.TwitsnapNotFound(postID)
		} else {
			return nil, postErrors.DatabaseError(err.Error())
		}
	}

	modPost, err = addAuthorInfoToPost(modPost, token)

	if err != nil {
		return nil, postErrors.UserInfoError(err.Error())
	}

	return &modPost, nil
}

func (c *Service) RetweetPost(postId string, userID string, token string) (*models.FrontPost, error) {
	post, err := c.db.GetPostByID(postId, userID)

	if err != nil {
		return nil, postErrors.TwitsnapNotFound(postId)
	}

	retweet := models.NewRetweetDBPost(post, userID)

	newRetweet, err := c.db.AddNewPost(retweet)

	if err != nil {
		return nil, postErrors.DatabaseError(err.Error())
	}

	newRetweet, err = addAuthorInfoToPost(newRetweet, token)

	if err != nil {
		return nil, postErrors.UserInfoError(err.Error())
	}

	return &newRetweet, nil
}

func (c *Service) FetchAllPosts(limitConfig models.LimitConfig, token string) ([]models.FrontPost, bool, error) {

	posts, hasMore, err := c.db.GetAllPosts(limitConfig, database.ADMIN)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	posts, err = addAuthorInfoToPosts(posts, token)

	return posts, hasMore, err
}

func (c *Service) FetchUserFeed(feedRequest *models.FeedRequesst, user_id string, limitConfig models.LimitConfig, token string) ([]models.FrontPost, bool, error) {
	switch feedRequest.FeedType {
	case FOLLOWING:
		return c.fetchFollowingFeed(limitConfig, user_id, token)
	case FORYOU:
		return c.fetchForyouFeed(limitConfig, user_id, token)
	case SINGLE:
		return c.fetchForyouSingle(limitConfig, feedRequest.WantedUserID, user_id, token)
	}
	return []models.FrontPost{}, false, postErrors.BadFeedRequest(feedRequest.FeedType)
}

func (c *Service) fetchFollowingFeed(limitConfig models.LimitConfig, userID string, token string) ([]models.FrontPost, bool, error) {
	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	log.Println("following: ", following)
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

	log.Println("interests: ", interests)

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

func (c *Service) FetchUserPostsByHashtags(hashtags []string, limitConfig models.LimitConfig, userID string, token string) ([]models.FrontPost, bool, error) {

	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}

	posts, hasMore, err := c.db.GetUserHashtags(hashtags, following, userID, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}

	posts, err = addAuthorInfoToPosts(posts, token)

	return posts, hasMore, err
}

func (c *Service) WordsSearch(words string, limitConfig models.LimitConfig, userID string, token string) ([]models.FrontPost, bool, error) {
	following, err := getUserFollowingWp(userID, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	posts, hasMore, err := c.db.WordSearchPosts(words, following, userID, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, nil
	}


	posts, err = addAuthorInfoToPosts(posts, token)

	return posts, hasMore, err
}

func (c *Service) LikePost(postID string, userID string) error {
	err := c.db.LikeAPost(postID, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}

func (c *Service) UnLikePost(postID string, userID string) error {
	err := c.db.UnLikeAPost(postID, userID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}
