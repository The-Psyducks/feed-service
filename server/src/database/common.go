package database

import (
	"server/src/models"
)

func allPostIntoFrontPost(posts []models.DBPost) []models.FrontPost {
	var frontPosts []models.FrontPost
	for _, post := range posts {
		frontPosts = append(frontPosts, makeDBPostIntoFrontPost(post))
	}
	return frontPosts
}

func makeDBPostIntoFrontPost(post models.DBPost) models.FrontPost {
	author := models.AuthorInfo{
		Author_ID: post.Author_ID,
		Username:  "username",
		Alias:     "alias",
		PthotoURL: "photourl",
	}
	return models.NewFrontPost(post, author)
}