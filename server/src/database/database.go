package database

import (
	"context"
	"log"
	"server/src/models"
	"strings"
	"time"

	postErrors "server/src/all_errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
)

type AppDatabase struct {
	db *mongo.Database
}

func NewAppDatabase(client *mongo.Client) Database {
	return &AppDatabase{db: client.Database(DATABASE_NAME)}
}

func (d *AppDatabase) AddNewPost(newPost models.DBPost) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	_, err := postCollection.InsertOne(context.Background(), newPost)

	frontPost := makeDBPostIntoFrontPost(newPost)
	return frontPost, err
}

func (d *AppDatabase) GetPostByID(postID string) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	post, err := d.findPost(postID, postCollection)

	frontPost := makeDBPostIntoFrontPost(post)
	return frontPost, err
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

func (d *AppDatabase) EditPost(postID string, editInfo models.EditPostExpectedFormat) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post models.FrontPost
	var dbPost models.DBPost

	err := d.updatePostContent(postID, editInfo.Content)

	if err != nil {
		return post, err
	}

	err_2 := d.updatePostTags(postID, editInfo.Tags)

	if err_2 != nil {
		return post, err_2
	}

	dbPost, err = d.findPost(postID, postCollection)

	frontPost := makeDBPostIntoFrontPost(dbPost)

	return frontPost, err
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

func (d *AppDatabase) GetUserFeedFollowing(following []string, limitConfig models.LimitConfig) ([]models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{AUTHOR_ID_FIELD: bson.M{"$in": following}, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}}

	cursor, err := postCollection.Find(context.Background(), filter,options.Find().
					SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)))
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	posts, err := createPostList(cursor, following)

	frontPosts := allPostIntoFrontPost(posts)

	return frontPosts, err
}

func (d *AppDatabase) GetUserFeedInterests(interests []string, following []string, limitConfig models.LimitConfig) ([]models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{TAGS_FIELD: bson.M{"$in": interests}, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}}

	cursor, err := postCollection.Find(context.Background(), filter, options.Find().
						SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)))
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	posts, err := createPostList(cursor, following)

	frontPosts := allPostIntoFrontPost(posts)

	return frontPosts, err
}

func (d *AppDatabase) GetUserFeedSingle(userId string, limitConfig models.LimitConfig) ([]models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var posts []models.DBPost

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{AUTHOR_ID_FIELD: userId, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}}

	cursor, err := postCollection.Find(context.Background(), filter, options.Find().
						SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)))
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

	frontPosts := allPostIntoFrontPost(posts)

	return frontPosts, err
}

func (d *AppDatabase) GetUserHashtags(interests []string, following []string, limitConfig models.LimitConfig) ([]models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)

	if len(interests) == 0 {
		return nil, postErrors.NoTagsFound()
	}

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{TAGS_FIELD: bson.M{"$all": interests}, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}}

	cursor, err := postCollection.Find(context.Background(), filter, options.Find().
						SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)))
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	posts, err := createPostList(cursor, following)

	frontPosts := allPostIntoFrontPost(posts)

	return frontPosts, err
}

func (d *AppDatabase) WordSearchPosts(words string, following []string, limitConfig models.LimitConfig) ([]models.FrontPost, error) {

	postCollection := d.db.Collection(FEED_COLLECTION)

	filters := bson.A{}

	for _, word := range strings.Split(words, " ") {
		if word != "" {
			log.Println(word)
			filters = append(filters, bson.M{CONTENT_FIELD: bson.M{"$regex": word, "$options": "i"}})
		}
	}

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{"$or": filters, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}}

	cursor, err := postCollection.Find(context.Background(), filter, options.Find().
						SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)))

	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	posts, err := createPostList(cursor, following)

	if posts == nil {
		err = postErrors.NoWordssFound()
	}

	frontPosts := allPostIntoFrontPost(posts)

	return frontPosts, err
}

func (d *AppDatabase) LikeAPost(postID string)  error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$inc": bson.M{LIKES_FIELD: 1}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) UnLikeAPost(postID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$inc": bson.M{LIKES_FIELD: -1}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}

func createPostList(cursor *mongo.Cursor, following []string) ([]models.DBPost, error) {
	var posts []models.DBPost
	var err error

	for cursor.Next(context.Background()) {
		var dbPost models.DBPost
		err = cursor.Decode(&dbPost)
		if err != nil {
			log.Println(err)
		}
		if dbPost.Public || (!dbPost.Public && slices.Contains(following, dbPost.Author_ID)) {
			posts = append(posts, dbPost)
		}
	}

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
