package service

import (
	"encoding/json"
	"errors"
	"io"
	// "log"
	"net/http"
	"os"
	"server/src/models"
	"strconv"
)

const (
	TEST_USER_ONE   = "1"
	TEST_USER_TWO   = "2"
	TEST_USER_THREE = "3"

	INITIAL_SKIP = 0

	TEST_TAG_ONE   = "tag1"
	TEST_TAG_TWO   = "tag2"
	TEST_TAG_THREE = "tag3"

	TEST_NOT_FOLLOWING_ID = "4"
)

func getUserFollowingWp(userID string, limitConfig models.LimitConfig, token string) ([]string, error) {
	if os.Getenv("ENVIROMENT") == "test" {
		return []string{TEST_USER_ONE, TEST_USER_TWO, TEST_USER_THREE}, nil
	} else {

		return getUserFollowing(userID, []string{}, limitConfig, INITIAL_SKIP, token)
	}
}

func getUserFollowing(userID string, following []string, limitConfig models.LimitConfig, skip int, token string) ([]string, error) {

	limit := strconv.Itoa(limitConfig.Limit)
	skipStr := strconv.Itoa(skip)

	// log.Println("userID: ", userID)
	// log.Println("token: ", token)

	url := "http://" + os.Getenv("USERS_HOST") + "/users/" + userID + "/following" + "?timestamp=" + limitConfig.FromTime + "&skip=" + skipStr + "&limit=" + limit

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		return nil, errors.New("error creating request")
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.New("error sending request, " + err.Error())
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	user := struct {
		Data       []models.UserInfoExpectedFormat `json:"data"`
		Pagination models.Pagination               `json:"pagination"`
	}{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		return nil, err
	}

	for _, data := range user.Data {
		following = append(following, data.Profile.ID)
	}

	if user.Pagination.Next_Offset != 0 {

		newLimit := models.NewLimitConfig(limitConfig.FromTime, limit, strconv.Itoa(user.Pagination.Next_Offset+limitConfig.Skip))

		return getUserFollowing(userID, following, newLimit, skip+limitConfig.Skip, token)
	}

	return following, nil
}

func getUserData(userID string, token string) (models.AuthorInfo, error) {

	url := "http://" + os.Getenv("USERS_HOST") + "/users/" + userID

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		return models.AuthorInfo{}, errors.New("error creating request")
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return models.AuthorInfo{}, errors.New("error sending request, " + err.Error())
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return models.AuthorInfo{}, errors.New("error reading request, " + err.Error())
	}

	user := struct {
		Following bool                               `json:"following"`
		Profile   models.PublicProfileExpectedFormat `json:"profile"`
	}{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		return models.AuthorInfo{}, errors.New("error unmarshaling request, " + err.Error())
	}

	authorInfo := models.AuthorInfo{Author_ID: user.Profile.ID, Username: user.Profile.Username,
		Alias: user.Profile.FisrtName + " " + user.Profile.LastName, PthotoURL: ""}

	return authorInfo, nil
}

func getUserInterestsWp(userID string, token string) ([]string, error) {
	if os.Getenv("ENVIROMENT") == "test" {
		return []string{TEST_TAG_ONE, TEST_TAG_TWO, TEST_TAG_THREE}, nil
	} else {
		return getUsersInterests(userID, token)
	}
}

func getUsersInterests(userID string, token string) ([]string, error) {
	url := "http://" + os.Getenv("USERS_HOST") + "/users/" + userID

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		return []string{}, errors.New("error creating request")
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return []string{}, errors.New("error sending request, " + err.Error())
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []string{}, errors.New("error reading request, " + err.Error())
	}

	user := struct {
		Following bool                                `json:"following"`
		Profile   models.PrivateProfileExpectedFormat `json:"profile"`
	}{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		return []string{}, errors.New("error unmarshaling request, " + err.Error())
	}

	return user.Profile.Interests, nil
}

func addAuthorInfoToPost(post models.FrontPost, token string) (models.FrontPost, error) {

	var authorInfo models.AuthorInfo
	var err error

	if os.Getenv("ENVIROMENT") == "test" {
		authorInfo, err = getUserDataForTests(post)
	} else {
		authorInfo, err = getUserData(post.Author_Info.Author_ID, token)
	}

	if err != nil {
		return models.FrontPost{}, errors.New("error getting info on the user, " + err.Error())
	}

	post.Author_Info = authorInfo

	if post.IsRetweet {
		post, err = addRetweetAuthorInfoToPost(post, token)
		if err != nil {
			return models.FrontPost{}, errors.New("error getting info on the user, " + err.Error())
		}
	} else {
		post.RetweetAuthor = ""
	}

	return post, nil
}

func addRetweetAuthorInfoToPost(post models.FrontPost, token string) (models.FrontPost, error) {
	var authorInfo models.AuthorInfo
	var err error

	if os.Getenv("ENVIROMENT") == "test" {
		authorInfo, err = getUserDataForTests(post)
	} else {
		authorInfo, err = getUserData(post.RetweetAuthor, token)
	}

	if err != nil {
		return models.FrontPost{}, errors.New("error getting info on the user, " + err.Error())
	}

	post.RetweetAuthor = authorInfo.Username

	return post, nil
}

func addAuthorInfoToPosts(posts []models.FrontPost, token string) ([]models.FrontPost, error) {
	for i, post := range posts {
		post, err := addAuthorInfoToPost(post, token)
		if err != nil {
			return nil, err
		}

		if post.IsRetweet {
			post, err = addRetweetAuthorInfoToPost(post, token)
			if err != nil {
				return []models.FrontPost{}, errors.New("error getting info on the user, " + err.Error())
			}
		} else {
			post.RetweetAuthor = ""
		}

		posts[i] = post
	}
	return posts, nil
}

func getUserDataForTests(post models.FrontPost) (models.AuthorInfo, error) {

	authorInfo := models.AuthorInfo{Author_ID: post.Author_Info.Author_ID, Username: "username",
		Alias: "alias", PthotoURL: ""}

	return authorInfo, nil
}
