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

func (c *Service) CreatePost(newPost *models.PostExpectedFormat) (*models.DBPost, error) {

	validate := validator.New()
	if err := validate.Struct(newPost); err != nil {
		return nil, postErrors.TwitSnapImportantFieldsMissing(err)
	}

	if len(newPost.Content) > 280 {
		return nil, postErrors.TwitsnapTooLong()
	}

	postNew := models.NewDBPost(newPost.Author_ID, newPost.Content, newPost.Tags, newPost.Public)

	if err := c.db.AddNewPost(postNew); err != nil {
		return nil, postErrors.DatabaseError()
	}

	return &postNew, nil
}

func (c *Service) FetchPostByID(postID string) (*models.DBPost, error) {

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

func (c *Service) ModifyPostByID(postID string, editInfo models.EditPostExpectedFormat) (*models.DBPost, error) {
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

func (c *Service) FetchUserFeed(userID string, feedType string) ([]models.DBPost, error) {
	switch feedType {
	case FOLLOWING:
		return c.fetchFollowingFeed(userID)
	case FORYOU:
		return c.fetchForyouFeed(userID)
	case SINGLE:
		return c.fetchForyouSingle(userID)
	}
	return nil, postErrors.BadFeedRequest()
}

func (c *Service) fetchFollowingFeed(userID string) ([]models.DBPost, error) {
	_ = userID
	following := []string{"3", "1"}
	posts, err := c.db.GetUserFeedFollowing(following)
	return posts, err
}

func (c *Service) fetchForyouFeed(userID string) ([]models.DBPost, error) {
	_ = userID
	following := []string{"apple", "1"}
	posts, err := c.db.GetUserFeedInterests(following)
	return posts, err
}

func (c *Service) fetchForyouSingle(userID string) ([]models.DBPost, error) {
	posts, err := c.db.GetUserFeedSingle(userID)
	return posts, err
}

func (c *Service) FetchUserPostsByHashtags(hashtags []string) ([]models.DBPost, error) {

	posts, err := c.db.GetUserHashtags(hashtags)

	if err != nil {
		return nil, err
	}

	if posts == nil {
		return nil, postErrors.NoTagsFound()
	}

	return posts, nil
}

func (c *Service) WordsSearch(words string) ([]models.DBPost, error) {
	posts, err := c.db.WordSearchPosts(words)

	if err != nil {
		return nil, err
	}

	if posts == nil {
		return nil, postErrors.NoWordssFound()
	}

	return posts, nil
}
