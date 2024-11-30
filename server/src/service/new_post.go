package service

import (
	"server/src/models"

	"github.com/go-playground/validator/v10"
	postErrors "server/src/all_errors"

	"strings"
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

	if c.amqp!= nil {
		if err := c.sendNewContentMessage(newPosted.Tags, newPosted.Time); err != nil {
			return nil, postErrors.QueueError(err.Error())
		}
	}

	return &newPosted, nil
}

func (c *Service) parsePost(post *models.PostExpectedFormat, author_id string) (models.DBPost, error) {
	validate := validator.New()
	if err := validate.Struct(post); err != nil {
		return models.DBPost{}, postErrors.TwitSnapImportantFieldsMissing(err)
	}

	if len(post.Content) > 280 {
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