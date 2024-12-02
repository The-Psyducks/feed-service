package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	postErrors "server/src/all_errors"
	"server/src/models"

)

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