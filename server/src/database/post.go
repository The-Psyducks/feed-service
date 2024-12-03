package database

import (
	"context"
	"log"
	postErrors "server/src/all_errors"
	"server/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *AppDatabase) AddNewPost(newPost models.DBPost) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	_, err := postCollection.InsertOne(context.Background(), newPost)

	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	frontPost, err_2 := d.makeDBPostIntoFrontPost(newPost, newPost.Author_ID)

	return frontPost, err_2
}

func (d *AppDatabase) GetPost(postID string, askerID string) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	post, err := d.findPost(postID, postCollection)

	if err != nil {
		return models.FrontPost{}, err
	}

	frontPost, err := d.makeDBPostIntoFrontPost(post, askerID)

	return frontPost, err
}

func (d *AppDatabase) DeletePost(postID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	likesCollection := d.db.Collection(LIKES_COLLECTION)
	retweetCollection := d.db.Collection(RETWEET_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}

	result, err := postCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return postErrors.ErrTwitsnapNotFound
	}

	filter_retweet := bson.M{ORIGINAL_POST_ID_FIELD: postID}
	_, err = postCollection.DeleteMany(context.Background(), filter_retweet)

	if err != nil {
		return err
	}

	_, err = likesCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	_, err = retweetCollection.DeleteOne(context.Background(), filter)

	return err
}

func (d *AppDatabase) GetAllPosts(limitConfig models.LimitConfig, askerID string) ([]models.FrontPost, bool, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}}

	cursor, err := postCollection.Find(context.Background(), filter, options.Find().
		SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)+1))

	if err != nil {
		log.Println(err)
		return []models.FrontPost{}, false, postErrors.DatabaseError(err.Error())
	}
	defer cursor.Close(context.Background())

	posts, err := d.createPostList(cursor, askerID)

	if err != nil {
		log.Println(err)
		return []models.FrontPost{}, false, postErrors.DatabaseError(err.Error())
	}

	hasMore := len(posts) > limitConfig.Limit

	if hasMore {
		posts = posts[:len(posts)-1]
	}

	return posts, hasMore, err
}