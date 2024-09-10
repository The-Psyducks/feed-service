package database

import (
	"log"
	"server/src/post"
	"sort"
	"strings"

	postErrors "server/src/all_errors"

	"github.com/mjarkk/mongomock"
	"go.mongodb.org/mongo-driver/bson"

	constants "server/src"
)

type TestDatabase struct {
	db *mongomock.TestConnection
}

func NewTestDatabase() Database {
	return &TestDatabase{db: mongomock.NewDB()}
}

func (d *TestDatabase) AddNewPost(newPost post.DBPost) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	err := postCollection.Insert(newPost)
	if err != nil {
		return err
	}
	return err
}

func (d *TestDatabase) GetPostByID(postID string) (post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	post, err := d.findPost(postID, postCollection)
	return post, err
}

func (d *TestDatabase) DeletePostByID(postID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{constants.POST_ID_FIELD: postID}

	err := postCollection.DeleteFirst(filter)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (d *TestDatabase) UpdatePostContent(postID string, newContent string) (post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post post.DBPost

	filter := bson.M{constants.POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{constants.CONTENT_FIELD: newContent}}

	err := postCollection.ReplaceFirst(filter, update)
	if err != nil {
		log.Println(err)
	}

	post, err = d.findPost(postID, postCollection)

	return post, err
}

func (d *TestDatabase) UpdatePostTags(postID string, newTags []string) (post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post post.DBPost

	filter := bson.M{constants.POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{constants.TAGS_FIELD: newTags}}

	err := postCollection.ReplaceFirst(filter, update)
	if err != nil {
		log.Println(err)
	}

	post, err = d.findPost(postID, postCollection)

	return post, err
}

func (d *TestDatabase) GetUserFeed(following []string) ([]post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []post.DBPost

	filter := bson.M{constants.AUTHOR_ID_FIELD: bson.M{"$in": following}}

	cursor, err := postCollection.FindCursor(filter)
	if err != nil {
		log.Println(err)
	}

	for cursor.Next() {
		var dbPost post.DBPost
		err := cursor.Decode(&dbPost)
		if err != nil {
			log.Println(err)
		}
		posts = append(posts, dbPost)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, err
}

func (d *TestDatabase) GetUserInterests(interests []string) ([]post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []post.DBPost

	filter := bson.M{constants.TAGS_FIELD: bson.M{"$all": interests}}

	cursor, err := postCollection.FindCursor(filter)
	if err != nil {
		log.Println(err)
	}

	for cursor.Next() {
		var dbPost post.DBPost
		err := cursor.Decode(&dbPost)
		if err != nil {
			log.Println(err)
		}
		if dbPost.Public {
			posts = append(posts, dbPost)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, err
}

func (d *TestDatabase) WordSearchPosts(words string) ([]post.DBPost, error) {
	
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []post.DBPost

	filters := bson.A{}

	for _, word := range strings.Split(words, " ") {
		log.Println(word)
		filters = append(filters, bson.M{constants.CONTENT_FIELD: bson.M{"$regex": word, "$options": "i"}})
	}

	filter := bson.M{"$or": filters}

	cursor, err := postCollection.FindCursor(filter)
	if err != nil {
		log.Println(err)
	}

	for cursor.Next() {
		var dbPost post.DBPost
		err := cursor.Decode(&dbPost)
		if err != nil {
			log.Println(err)
		}
		if dbPost.Public {
			posts = append(posts, dbPost)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, err
}

func (d *TestDatabase) findPost(postID string, postCollection *mongomock.Collection) (post.DBPost, error) {
	var post post.DBPost
	filter := bson.M{constants.POST_ID_FIELD: postID}
	err := postCollection.FindFirst(&post, filter)
	if err != nil {
		log.Println(err)
		err = postErrors.ErrTwitsnapNotFound
	}
	return post, err
}