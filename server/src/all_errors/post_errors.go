package allerrors

import (
	"errors"
	"net/http"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type TwitsnapError interface {
	Status() int
	Error() string
}

type TwitSnapError struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	ErrorStatus int    `json:"status"`
	Detail      string `json:"detail"`
	Instance    string `json:"instance"`
}

func (e TwitSnapError) Error() string {
	return e.Detail
}

func (e TwitSnapError) Status() int {
	return e.ErrorStatus
}

var ErrTwitsnapNotFound = errors.New("twitsnap not found")

func TwitsnapNotFound(id string) TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Twitsnap Not Found",
		http.StatusNotFound,
		"The twitsnap with ID " + id + " was not found",
		"/twitsnap/" + id,
	}
	return error
}

func TwitsnapTooLong() TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Twitsnap Too Long",
		http.StatusRequestEntityTooLarge,
		"The twitsnap message is too long",
		"/twitsnap",
	}
	return error
}

func TwitSnapImportantFieldsMissing(err error) TwitSnapError {
	missingFields := []string{}

	if validationError, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationError {
			missingFields = append(missingFields, fieldError.Field())
		}
	}

	error := TwitSnapError{
		"about:blank",
		"Twitsnap Important Fields Missing",
		http.StatusBadRequest,
		"The twitsnap message is missing " + strings.Join(missingFields, ","),
		"/twitsnap",
	}
	return error
}

func UnexpectedFormat() TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Unexpected Format",
		http.StatusBadRequest,
		"The twitsnap message has an unexpected format",
		"/twitsnap",
	}
	return error
}

func BadFeedRequest(feed string) TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Unexpected Format",
		http.StatusBadRequest,
		"There is no feed like that: " + feed,
		"/twitsnap",
	}
	return error
}

func NoTagsFound() TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"No posts with the tags found",
		http.StatusBadRequest,
		"There are no twitsnaps with the tags provided",
		"/twitsnap",
	}
	return error
}

func NoWordssFound() TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"No posts with the tags found",
		http.StatusBadRequest,
		"There are no twitsnaps with any of the words provided",
		"/twitsnap",
	}
	return error
}

func UserInfoError(err string) TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Error getting user info",
		http.StatusInternalServerError,
		err,
		"/twitsnap",
	}
	return error
}

func AuthenticationErrorHeaderRequired() TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Authorization header is required",
		http.StatusUnauthorized,
		"authorization header is required",
		"/twitsnap",
	}
	return error
}

func AuthenticationErrorInvalidHeader() TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Invalid authorization header",
		http.StatusUnauthorized,
		"Invalid authorization header",
		"/twitsnap",
	}
	return error
}

func AuthenticationErrorInvalidToken(err string) TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Invalid token",
		http.StatusUnauthorized,
		err,
		"/twitsnap",
	}
	return error
}

func DatabaseError(err string) TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Database Error",
		http.StatusInternalServerError,
		"There was an error with the database: " + err,
		"/twitsnap",
	}
	return error
}

func AccssDenied() TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Access Denied",
		http.StatusForbidden,
		"Access Denied, user not Admin",
		"/twitsnap",
	}
	return error
}

func NotFound() TwitSnapError {
	error := TwitSnapError{
		"about:blank",
		"Not Found",
		http.StatusNotFound,
		"Resource not found",
		"/twitsnap",
	}
	return error
}