package service

import (
	"errors"
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

func (c *Service) CreatePost(newPost *models.PostExpectedFormat) (*models.FrontPost, error) {

	validate := validator.New()
	if err := validate.Struct(newPost); err != nil {
		return nil, postErrors.TwitSnapImportantFieldsMissing(err)
	}

	if len(newPost.Content) > 280 {
		return nil, postErrors.TwitsnapTooLong()
	}

	postNew := models.NewDBPost(newPost.Author_ID, newPost.Content, newPost.Tags, newPost.Public)

	newPosted, err := c.db.AddNewPost(postNew)

	if err != nil {
		return nil, postErrors.DatabaseError()
	}

	return &newPosted, nil
}

func (c *Service) FetchPostByID(postID string) (*models.FrontPost, error) {

	post, err := c.db.GetPostByID(postID)

	if err != nil {
		return nil, postErrors.TwitsnapNotFound(postID)
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

func (c *Service) ModifyPostByID(postID string, editInfo models.EditPostExpectedFormat) (*models.FrontPost, error) {
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

	return &modPost, nil
}

func (c *Service) FetchUserFeed(userID string, feedType string, limitConfig models.LimitConfig) (models.ReturnPaginatedPosts, error) {
	switch feedType {
	case FOLLOWING:
		return c.fetchFollowingFeed(userID, limitConfig)
	case FORYOU:
		return c.fetchForyouFeed(userID, limitConfig)
	case SINGLE:
		return c.fetchForyouSingle(userID, limitConfig)
	}
	return models.ReturnPaginatedPosts{}, postErrors.BadFeedRequest()
}

func (c *Service) fetchFollowingFeed(userID string, limitConfig models.LimitConfig) (models.ReturnPaginatedPosts, error) {
	_ = userID
	following := []string{"3", "1"}
	posts, err := c.db.GetUserFeedFollowing(following, limitConfig)
	return posts, err
}

func (c *Service) fetchForyouFeed(userID string, limitConfig models.LimitConfig) (models.ReturnPaginatedPosts, error) {
	_ = userID
	interests := []string{"apple", "1"}
	following := []string{"3", "1"}
	posts, err := c.db.GetUserFeedInterests(interests, following, limitConfig)
	return posts, err
}

func (c *Service) fetchForyouSingle(userID string, limitConfig models.LimitConfig) (models.ReturnPaginatedPosts, error) {
	posts, err := c.db.GetUserFeedSingle(userID, limitConfig)
	return posts, err
}

func (c *Service) FetchUserPostsByHashtags(hashtags []string, limitConfig models.LimitConfig) (models.ReturnPaginatedPosts, error) {
	following := []string{"3", "1"}

	posts, err := c.db.GetUserHashtags(hashtags, following, limitConfig)

	if err != nil {
		return models.ReturnPaginatedPosts{}, err
	}

	if posts.Data == nil {
		return models.ReturnPaginatedPosts{}, postErrors.NoTagsFound()
	}

	return posts, nil
}

func (c *Service) WordsSearch(words string, limitConfig models.LimitConfig) (models.ReturnPaginatedPosts, error) {
	following := []string{"3", "1"}
	posts, err := c.db.WordSearchPosts(words, following, limitConfig)

	if err != nil {
		return models.ReturnPaginatedPosts{}, err
	}

	if posts.Data == nil {
		return models.ReturnPaginatedPosts{}, postErrors.NoWordssFound()
	}

	return posts, nil
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
