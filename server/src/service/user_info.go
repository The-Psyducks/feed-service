package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"server/src/models"
	"strconv"
)

func getUserFollowingWp(username string, limitConfig models.LimitConfig, token string) ([]string, error) {
	return getUserFollowing(username, []string{}, limitConfig, token)

}

func getUserFollowing(username string, following []string, limitConfig models.LimitConfig, token string) ([]string, error) {

	limit := strconv.Itoa(limitConfig.Limit)
	skip := strconv.Itoa(limitConfig.Skip)

	url := "http://" + os.Getenv("USERS_HOST") + "/users/" + username + "/following" + "?timestamp=" + limitConfig.FromTime + "&skip=" + skip + "&limit=" + limit

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

		return getUserFollowing(username, following, newLimit, token)
	}

	return following, nil
}

func getUserData(userID string, token string) (models.AuthorInfo, error) {

	url := "http://" + os.Getenv("USERS_HOST") + "/users/profile/" + userID

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

func getUsernameData(username string, token string) (models.AuthorInfo, error) {
	url := "http://" + os.Getenv("USERS_HOST") + "/users/" + username

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

func addAuthorInfoToPost(post models.FrontPost, token string) (models.FrontPost, error) {

	authorInfo, err := getUserData(post.Author_Info.Author_ID, token)

	if err != nil {
		return models.FrontPost{}, errors.New("error getting info on the user, " + err.Error())
	}

	post.Author_Info = authorInfo

	return post, nil
}

func addAuthorInfoToPosts(posts []models.FrontPost, token string) ([]models.FrontPost, error) {

	for i, post := range posts {
		post, err := addAuthorInfoToPost(post, token)

		if err != nil {
			return nil, err
		}

		posts[i] = post
	}

	return posts, nil
}

func getUserID(username string, token string) (string, error) {

	userData, err := getUsernameData(username, token)

	return userData.Author_ID, err
}