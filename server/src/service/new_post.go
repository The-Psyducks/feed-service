package service

import (
	"log/slog"
	"server/src/database"
	"server/src/models"
	"time"

	postErrors "server/src/all_errors"

	"github.com/go-playground/validator/v10"

	"strings"
)

const (
	MAX_CHAR = 280
)

func (c *Service) CreatePost(newPost *models.PostExpectedFormat, author_id string, token string) (*models.FrontPost, error) {

	postNew, err := c.parsePost(newPost, author_id)

	if err != nil {
		return nil, err
	}

	newPosted, err := c.db.AddNewPost(postNew)

	if err != nil {
		return nil, postErrors.DatabaseError(err.Error())
	}

	newPosted, err = addAuthorInfoToPost(newPosted, token)

	if err != nil {
		return nil, postErrors.UserInfoError(err.Error())
	}

	for _, user := range newPost.Mentions {
		newMentionNotif := models.MentionNotificationRequest{UserId: user, TaggerId: newPosted.Author_Info.Author_ID, PostId: newPosted.Original_Post_ID}
		err = sendMentionNotif(newMentionNotif, token)

		if err != nil {
			return nil, postErrors.NotificationError(err.Error())
		}
	}

	slog.Info("New post created: ", "original_post_id", newPosted.Original_Post_ID, "author_id", newPosted.Author_Info.Author_ID, "timestamp", newPosted.Time)

	return &newPosted, nil
}

func (c *Service) parsePost(post *models.PostExpectedFormat, author_id string) (models.DBPost, error) {
	validate := validator.New()
	if err := validate.Struct(post); err != nil {
		return models.DBPost{}, postErrors.TwitSnapImportantFieldsMissing(err)
	}

	if len(post.Content) > MAX_CHAR {
		return models.DBPost{}, postErrors.TwitsnapTooLong()
	}

	var tags []string

	content := strings.Split(post.Content, " ")

	for _, word := range content {
		if strings.HasPrefix(word, "#") {
			word = word[1:]
			tags = append(tags, word)
		}
	}

	postNew := models.NewDBPost(author_id, post.Content, tags, post.Public, post.MediaInfo, post.Mentions)

	return postNew, nil
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

	slog.Info("All posts retrieved: ", "time", time.Now(), "count", len(posts))

	return posts, hasMore, err
}