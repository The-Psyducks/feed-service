package database

import (
	"log"
	"server/src/models"
	"sort"
	"strings"

	postErrors "server/src/all_errors"

	"github.com/mjarkk/mongomock"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

type TestDatabase struct {
	db *mongomock.TestConnection
}

func NewTestDatabase() Database {
	return &TestDatabase{db: mongomock.NewDB()}
}

func (d *TestDatabase) AddNewPost(newPost models.DBPost) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	err := postCollection.Insert(newPost)
	
	frontPost := makeDBPostIntoFrontPost(newPost)

	return frontPost, err
}

func (d *TestDatabase) GetPostByID(postID string) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	post, err := d.findPost(postID, postCollection)

	frontPost := makeDBPostIntoFrontPost(post)
	return frontPost, err
}

func (d *TestDatabase) DeletePostByID(postID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}

	err := postCollection.DeleteFirst(filter)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (d *TestDatabase) EditPost(postID string, editInfo models.EditPostExpectedFormat) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post models.FrontPost
	var dbPost models.DBPost

	err := d.updatePostContent(postID, editInfo.Content)

	if err != nil {
		return post, err
	}

	err = d.updatePostTags(postID, editInfo.Tags)

	if err != nil {
		return post, err
	}

	dbPost, err = d.findPost(postID, postCollection)

	frontPost := makeDBPostIntoFrontPost(dbPost)

	return frontPost, err

}

func (d *TestDatabase) updatePostContent(postID string, newContent string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	if len(newContent) == 0 {
		return nil
	}

	post, err := d.findPost(postID, postCollection)

	if err != nil {
		return err
	}

	err_2 := d.DeletePostByID(postID)

	if err_2 != nil {
		return err_2
	}

	post.Content = newContent

	_ , err_3 := d.AddNewPost(post)

	return err_3
}

func (d *TestDatabase) updatePostTags(postID string, newTags []string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	if len(newTags) == 0 {
		return nil
	}

	post, err := d.findPost(postID, postCollection)

	if err != nil {
		return err
	}

	err_2 := d.DeletePostByID(postID)

	if err_2 != nil {
		return err_2
	}

	post.Tags = newTags

	_, err_3 := d.AddNewPost(post)

	return err_3
}

func (d *TestDatabase) GetUserFeedFollowing(following []string) ([]models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	cursor, err := postCollection.FindCursor(bson.M{})
	if err != nil {
		log.Println(err)
	}

	for cursor.Next() {
		var dbPost models.DBPost
		err := cursor.Decode(&dbPost)
		if err != nil {
			log.Println(err)
		}
		posts = append(posts, dbPost)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time > posts[j].Time
	})

	feed := []models.DBPost{}

	for _, post := range posts {
		if slices.Contains(following, post.Author_ID) {
			feed = append(feed, post)
		}
	}

	return allPostIntoFrontPost(feed), err
}

func (d *TestDatabase) GetUserHashtags(interests []string, following []string) ([]models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	cursor, err := postCollection.FindCursor(bson.M{})
	if err != nil {
		log.Println(err)
	}

	for cursor.Next() {
		var dbPost models.DBPost
		err := cursor.Decode(&dbPost)
		if err != nil {
			log.Println(err)
		}
		if dbPost.Public {
			posts = append(posts, dbPost)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time > posts[j].Time
	})

	feed := []models.DBPost{}

	for _, post := range posts {
		if containsAll(post.Tags, interests) {
			feed = append(feed, post)
		}
	}

	return allPostIntoFrontPost(feed), err
}

func containsAll(s []string, e []string) bool {
	for _, a := range e {
		if !slices.Contains(s, a) {
			return false
		}
	}
	return true
}

func (d *TestDatabase) WordSearchPosts(words string, following []string) ([]models.FrontPost, error) {

	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	filter := []string{}

	filter = append(filter, strings.Split(words, " ")...)

	cursor, err := postCollection.FindCursor(bson.M{})
	if err != nil {
		log.Println(err)
	}

	for cursor.Next() {
		var dbPost models.DBPost
		err := cursor.Decode(&dbPost)
		if err != nil {
			log.Println(err)
		}
		if dbPost.Public {
			posts = append(posts, dbPost)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time > posts[j].Time
	})

	result := []models.DBPost{}

	for _, post := range posts {
		postContent := strings.Split(post.Content, " ")

		if containsOne(postContent, filter) {
			result = append(result, post)
		}
	}

	return allPostIntoFrontPost(result), err
}

func containsOne(s []string, e []string) bool {
	for _, a := range e {
		if slices.Contains(s, a) {
			return true
		}
	}
	return false
}

func (d *TestDatabase) findPost(postID string, postCollection *mongomock.Collection) (models.DBPost, error) {
	var post models.DBPost
	filter := bson.M{POST_ID_FIELD: postID}
	err := postCollection.FindFirst(&post, filter)
	if err != nil {
		log.Println(err)
		err = postErrors.ErrTwitsnapNotFound
	}
	return post, err
}

func (d *TestDatabase) GetUserFeedInterests(interests []string, following []string) ([]models.FrontPost, error) {
	return nil, nil
}

func (d *TestDatabase) GetUserFeedSingle(userId string) ([]models.FrontPost, error) {
	return nil, nil
}
