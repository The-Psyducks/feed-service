package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	postErrors "server/src/all_errors"
	"server/src/database"
	"server/src/models"
	"strconv"

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

func (c *Service) FetchUserFeed(username string, feedType string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	switch feedType {
	case FOLLOWING:
		return c.fetchFollowingFeed(username, limitConfig)
	case FORYOU:
		return c.fetchForyouFeed(username, limitConfig)
	case SINGLE:
		return c.fetchForyouSingle(username, limitConfig)
	}
	return []models.FrontPost{}, false, postErrors.BadFeedRequest()
}

func (c *Service) fetchFollowingFeed(username string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	_ = username
	following := []string{"3", "1"}
	posts, hasMore, err := c.db.GetUserFeedFollowing(following, limitConfig)
	return posts, hasMore, err
}

func (c *Service) fetchForyouFeed(username string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	_ = username
	interests := []string{"apple", "1"}
	following := []string{"3", "1"}
	posts, hasMore, err := c.db.GetUserFeedInterests(interests, following, limitConfig)
	return posts, hasMore, err
}

func (c *Service) fetchForyouSingle(userID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	posts, hasMore, err := c.db.GetUserFeedSingle(userID, limitConfig)
	return posts, hasMore, err
}

func (c *Service) FetchUserPostsByHashtags(hashtags []string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	following := []string{"3", "1"}

	posts, hasMore, err := c.db.GetUserHashtags(hashtags, following, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, postErrors.NoTagsFound()
	}

	return posts, hasMore, nil
}

func (c *Service) WordsSearch(words string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	following := []string{"3", "1"}
	posts, hasMore, err := c.db.WordSearchPosts(words, following, limitConfig)

	if err != nil {
		return []models.FrontPost{}, false, err
	}

	if len(posts) == 0 {
		return []models.FrontPost{}, false, postErrors.NoWordssFound()
	}

	return posts, hasMore, nil
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

func getUserFollowingWp(username string, limitConfig models.LimitConfig) []string {
	return getUserFollowing(username, []string{}, limitConfig)
	
}

func getUserFollowing(username string, following []string, limitConfig models.LimitConfig) []string {

	limit := strconv.Itoa(limitConfig.Limit)
	skip := strconv.Itoa(limitConfig.Skip)
	
	url := "http://localhost:8080/users/" + username + "?time=" + limitConfig.FromTime + "&skip="+ skip +"&limit=" + limit

	req, err := http.Get(url)

	if err != nil {
		return following
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return following
	}

	user := struct {
		Data []models.UserFollowingExpectedFormat `json:"data"`
		Pagination models.Pagination `json:"pagination"`
	}{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		return following
	}

	for _, data := range user.Data {
		following = append(following, data.Profile.ID)
	}



	if user.Pagination.Next_Offset != 0 {

		newLimit := models.NewLimitConfig(limitConfig.FromTime, limit, strconv.Itoa(user.Pagination.Next_Offset + limitConfig.Skip))

		return getUserFollowing(username, following, newLimit)
	}

	return following
}
