package database

import (
	"context"
	"log"
	"server/src/post"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	postErrors "server/src/all_errors"

	constants "server/src"
)

const (
	DATABASE_NAME = "feed"
	FEED_COLLECTION = "posts"
)

type Database struct {
	db *mongo.Database
}

func NewDatabase(client *mongo.Client) *Database {
	return &Database{db: client.Database(DATABASE_NAME)}
}

func (d *Database) AddNewPost(newPost post.DBPost) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	_, err := postCollection.InsertOne(context.Background(), newPost)
	if err != nil {
		return err
	}
	return err
}

func (d *Database) GetPostByID(postID string) (post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	post, err := d.findPost(postID, postCollection)
	return post, err
}

func (d *Database) DeletePostByID(postID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{constants.POST_ID_FIELD: postID}

	_, err := postCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (d *Database) UpdatePostContent(postID string, newContent string) (post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post post.DBPost

	filter := bson.M{constants.POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{constants.CONTENT_FIELD: newContent}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	post, err = d.findPost(postID, postCollection)

	return post, err
}

func (d *Database) UpdatePosTags(postID string, newTags []string) (post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post post.DBPost

	filter := bson.M{constants.POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{constants.TAGS_FIELD: newTags}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	post, err = d.findPost(postID, postCollection)

	return post, err
}

func (d *Database) GetUserFeed(following []string) ([]post.DBPost, error) {
	return d.getFiltered(following, constants.AUTHOR_ID_FIELD, "$in")
}

func (d *Database) GetUserInterests(interests []string) ([]post.DBPost, error) {
	return d.getFiltered(interests, constants.TAGS_FIELD, "$all")
}

func (d *Database) getFiltered(dataFilter []string, field string, filterLogic string) ([]post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []post.DBPost

	filter := bson.M{field: bson.M{filterLogic: dataFilter}}

	cursor, err := postCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
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

func (d *Database) findPost(postID string, postCollection *mongo.Collection) (post.DBPost, error) {
	var post post.DBPost
	filter := bson.M{constants.POST_ID_FIELD: postID}
	err := postCollection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		log.Println(err)
		err = postErrors.ErrTwitsnapNotFound
	}
	return post, err
}