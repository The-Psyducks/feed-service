package database

import (
	"context"
	"log"
	"server/src/post"
	"sort"
	"strings"

	postErrors "server/src/all_errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	constants "server/src"
)

const (
	DATABASE_NAME = "feed"
	FEED_COLLECTION = "posts"
)

type AppDatabase struct {
	db *mongo.Database
}

func NewAppDatabase(client *mongo.Client) Database {
	return &AppDatabase{db: client.Database(DATABASE_NAME)}
}

func (d *AppDatabase) AddNewPost(newPost post.DBPost) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	_, err := postCollection.InsertOne(context.Background(), newPost)
	if err != nil {
		return err
	}
	return err
}

func (d *AppDatabase) GetPostByID(postID string) (post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	post, err := d.findPost(postID, postCollection)
	return post, err
}

func (d *AppDatabase) DeletePostByID(postID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{constants.POST_ID_FIELD: postID}

	_, err := postCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (d *AppDatabase) EditPost(postID string, editInfo post.EditPostExpectedFormat) (post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post post.DBPost

	err := d.updatePostContent(postID, editInfo.Content)

	if err != nil {
		return post, err
	}

	err_2 := d.updatePostTags(postID, editInfo.Tags)

	if err_2 != nil {
		return post, err_2
	}

	post, err = d.findPost(postID, postCollection)

	return post, err
}

func (d *AppDatabase) updatePostContent(postID string, newContent string)  error {

	if len(newContent) == 0 {
		return nil
	}

	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{constants.POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{constants.CONTENT_FIELD: newContent}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) updatePostTags(postID string, newTags []string) error {

	if len(newTags) == 0 {
		return nil
	}

	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{constants.POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{constants.TAGS_FIELD: newTags}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) GetUserFeed(following []string) ([]post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []post.DBPost

	filter := bson.M{constants.AUTHOR_ID_FIELD: bson.M{"$in": following}}

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

func (d *AppDatabase) GetUserInterests(interests []string) ([]post.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []post.DBPost

	filter := bson.M{constants.TAGS_FIELD: bson.M{"$all": interests}}

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
		if dbPost.Public {
			posts = append(posts, dbPost)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, err
}

func (d *AppDatabase) WordSearchPosts(words string) ([]post.DBPost, error) {
	
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []post.DBPost

	filters := bson.A{}

	for _, word := range strings.Split(words, " ") {
		log.Println(word)
		filters = append(filters, bson.M{constants.CONTENT_FIELD: bson.M{"$regex": word, "$options": "i"}})
	}

	filter := bson.M{"$or": filters}

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
		if dbPost.Public {
			posts = append(posts, dbPost)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, err
}

func (d *AppDatabase) findPost(postID string, postCollection *mongo.Collection) (post.DBPost, error) {
	var post post.DBPost
	filter := bson.M{constants.POST_ID_FIELD: postID}
	err := postCollection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		log.Println(err)
		err = postErrors.ErrTwitsnapNotFound
	}
	return post, err
}