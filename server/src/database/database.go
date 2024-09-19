package database

import (
	"context"
	"log"
	"server/src/models"
	"sort"
	"strings"

	postErrors "server/src/all_errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppDatabase struct {
	db *mongo.Database
}

func NewAppDatabase(client *mongo.Client) Database {
	return &AppDatabase{db: client.Database(DATABASE_NAME)}
}

func (d *AppDatabase) AddNewPost(newPost models.DBPost) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	_, err := postCollection.InsertOne(context.Background(), newPost)
	if err != nil {
		return err
	}
	return err
}

func (d *AppDatabase) GetPostByID(postID string) (models.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	post, err := d.findPost(postID, postCollection)
	return post, err
}

func (d *AppDatabase) DeletePostByID(postID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}

	result, err := postCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}

	if result.DeletedCount == 0 {
		err = postErrors.ErrTwitsnapNotFound
	}

	return err
}

func (d *AppDatabase) EditPost(postID string, editInfo models.EditPostExpectedFormat) (models.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post models.DBPost

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

func (d *AppDatabase) updatePostContent(postID string, newContent string) error {

	if len(newContent) == 0 {
		return nil
	}

	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{CONTENT_FIELD: newContent}}

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

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{TAGS_FIELD: newTags}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) GetUserFeedFollowing(following []string) ([]models.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	filter := bson.M{AUTHOR_ID_FIELD: bson.M{"$in": following}}

	cursor, err := postCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var dbPost models.DBPost
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

func (d *AppDatabase) GetUserFeedInterests(interests []string) ([]models.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	filter := bson.M{TAGS_FIELD: bson.M{"$in": interests}}

	cursor, err := postCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var dbPost models.DBPost
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

func (d *AppDatabase) GetUserFeedSingle(userId string) ([]models.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	filter := bson.M{AUTHOR_ID_FIELD: userId}

	cursor, err := postCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var dbPost models.DBPost
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

func (d *AppDatabase) GetUserHashtags(interests []string) ([]models.DBPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	filter := bson.M{TAGS_FIELD: bson.M{"$all": interests}}

	cursor, err := postCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
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
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, err
}

func (d *AppDatabase) WordSearchPosts(words string) ([]models.DBPost, error) {

	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	filters := bson.A{}

	for _, word := range strings.Split(words, " ") {
		if word != "" {
			log.Println(word)
			filters = append(filters, bson.M{CONTENT_FIELD: bson.M{"$regex": word, "$options": "i"}})
		}
	}

	filter := bson.M{"$or": filters}

	cursor, err := postCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var dbPost models.DBPost
		err := cursor.Decode(&dbPost)
		if err != nil {
			log.Println(err)
		}
		if dbPost.Public {
			posts = append(posts, dbPost)
		}
	}

	if posts == nil {
		err = postErrors.NoWordssFound()
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, err
}

func (d *AppDatabase) findPost(postID string, postCollection *mongo.Collection) (models.DBPost, error) {
	var post models.DBPost
	filter := bson.M{POST_ID_FIELD: postID}
	err := postCollection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		log.Println(err)
		err = postErrors.ErrTwitsnapNotFound
	}
	return post, err
}
