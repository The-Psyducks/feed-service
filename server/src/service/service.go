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

	postNew := models.NewDBPost(author_id, newPost.Content, newPost.Tags, newPost.Public)

	newPosted, err := c.db.AddNewPost(postNew)

	if err != nil {
		return nil, postErrors.DatabaseError()
	}

	newPosted, err = addAuthorInfoToPost(newPosted, token)

	if err != nil {
		return nil, postErrors.UserInfoError(err.Error())
	}

	return &newPosted, nil
}

func (c *Service) FetchPostByID(postID string, token string) (*models.FrontPost, error) {

	post, err := c.db.GetPostByID(postID)

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

func (c *Service) ModifyPostByID(postID string, editInfo models.EditPostExpectedFormat, token string) (*models.FrontPost, error) {
	validate := validator.New()
	if err := validate.Struct(editInfo); err != nil {
		return nil, postErrors.TwitSnapImportantFieldsMissing(err)
	}

	modPost, err := c.db.EditPost(postID, editInfo)

	if err != nil {
		if errors.Is(err, postErrors.ErrTwitsnapNotFound) {
			return nil, postErrors.TwitsnapNotFound(postID)
		} else {
			return nil, postErrors.DatabaseError()
		}
	}

	modPost, err = addAuthorInfoToPost(modPost, token)

	if err != nil {
		return nil, postErrors.UserInfoError(err.Error())
	}

	return &modPost, nil
}

func (c *Service) FetchUserFeed(feedRequest models.FeedRequesst, username string, limitConfig models.LimitConfig, token string) ([]models.FrontPost, bool, error) {

	validate := validator.New()
	if err := validate.Struct(feedRequest); err != nil {
		return nil, false, postErrors.TwitSnapImportantFieldsMissing(err)
	}

	switch feedRequest.FeedType {
	case FOLLOWING:
		return c.fetchFollowingFeed(username, limitConfig, token)
	case FORYOU:
		return c.fetchForyouFeed(username, limitConfig, token)
	case SINGLE:
		return c.fetchForyouSingle(limitConfig, username, token)
	}
	return []models.FrontPost{}, false, postErrors.BadFeedRequest()
}

func (c *Service) fetchFollowingFeed(username string, limitConfig models.LimitConfig, token string) ([]models.FrontPost, bool, error) {
	_ = username
	following, err := getUserFollowingWp(username, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	log.Println(following)
	posts, hasMore, err := c.db.GetUserFeedFollowing(following, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	posts, err = addAuthorInfoToPosts(posts, token)
	return posts, hasMore, err
}

func (c *Service) fetchForyouFeed(username string, limitConfig models.LimitConfig, token string) ([]models.FrontPost, bool, error) {
	_ = username
	interests := []string{"apple", "1"}
	following, err := getUserFollowingWp(username, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	posts, hasMore, err := c.db.GetUserFeedInterests(interests, following, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	posts, err = addAuthorInfoToPosts(posts, token)
	return posts, hasMore, err
}

func (c *Service) fetchForyouSingle(limitConfig models.LimitConfig, username string, token string) ([]models.FrontPost, bool, error) {
	userID, err := getUserID(username, token)
	if err != nil {
		return []models.FrontPost{}, false, postErrors.UserInfoError(err.Error())
	}
	following, err := getUserFollowingWp(username, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	posts, hasMore, err := c.db.GetUserFeedSingle(userID, limitConfig, following)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	posts, err = addAuthorInfoToPosts(posts, token)
	return posts, hasMore, err
}

func (c *Service) FetchUserPostsByHashtags(hashtags []string, limitConfig models.LimitConfig, username string, token string) ([]models.FrontPost, bool, error) {
	following, err := getUserFollowingWp(username, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}

	posts, hasMore, err := c.db.GetUserHashtags(hashtags, following, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, postErrors.NoTagsFound()
	}

	posts, err = addAuthorInfoToPosts(posts, token)

	return posts, hasMore, err
}

func (c *Service) WordsSearch(words string, limitConfig models.LimitConfig, username string, token string) ([]models.FrontPost, bool, error) {
	following, err := getUserFollowingWp(username, limitConfig, token)
	if err != nil {
		return []models.FrontPost{}, false, err
	}
	posts, hasMore, err := c.db.WordSearchPosts(words, following, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, postErrors.NoWordssFound()
	}

	posts, err = addAuthorInfoToPosts(posts, token)

	return posts, hasMore, err
}

func (c *Service) LikePost(postID string) error {
	err := c.db.LikeAPost(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}

func (c *Service) UnLikePost(postID string) error {
	err := c.db.UnLikeAPost(postID)

	if err != nil {
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}

