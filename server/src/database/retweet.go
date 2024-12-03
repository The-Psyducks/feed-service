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

func (d *AppDatabase) AddNewRetweet(newRetweet models.DBPost) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	_, err := postCollection.InsertOne(context.Background(), newRetweet)

	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	retweetCollection := d.db.Collection(RETWEET_COLLECTION)

	filter_original := bson.M{ORIGINAL_POST_ID_FIELD: newRetweet.Original_Post_ID}
	update := bson.M{"$inc": bson.M{RETWEET_FIELD: 1}}

	retweeter := bson.M{"$addToSet": bson.M{RETWEETERS_FIELD: newRetweet.Retweet_Author_ID}}

	_, err = postCollection.UpdateMany(context.Background(), filter_original, update)
	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	_, err = retweetCollection.UpdateOne(context.Background(), filter_original, retweeter, options.Update().SetUpsert(true))
	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	newPostRetweet, err := d.findPost(newRetweet.Post_ID, postCollection)

	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	post, err := d.makeDBPostIntoFrontPost(newPostRetweet, newPostRetweet.Retweet_Author_ID)

	return post, err
}

func (d *AppDatabase) DeleteRetweet(postID string, userID string) error {

	postCollection := d.db.Collection(FEED_COLLECTION)
	retweetCollection := d.db.Collection(RETWEET_COLLECTION)

	filter := bson.M{ORIGINAL_POST_ID_FIELD: postID}
	update := bson.M{"$inc": bson.M{RETWEET_FIELD: -1}}

	retweeter := bson.M{"$pull": bson.M{RETWEETERS_FIELD: userID}}

	_, err := postCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = retweetCollection.UpdateOne(context.Background(), filter, retweeter)

	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) GetUserFeedRetweet(userId string, limitConfig models.LimitConfig, askerID string, following []string) ([]models.FrontPost, bool, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
		log.Println("User does not follow")
	}

	filter := bson.M{
		TIME_FIELD:       bson.M{"$lt": parsedTime.UTC()},
		IS_RETWEET_FIELD: true,
		"$and": []bson.M{
			{"$or": []bson.M{
				{AUTHOR_ID_FIELD: userId},
				{RETWEET_AUTHOR_FIELD: userId},
			}},
			{"$or": []bson.M{
				{PUBLIC_FIELD: true},
				{PUBLIC_FIELD: false, AUTHOR_ID_FIELD: bson.M{"$in": following}},
				{PUBLIC_FIELD: false, RETWEET_AUTHOR_FIELD: bson.M{"$in": following}},
			}},
		},
	}

	cursor, err := postCollection.Find(context.Background(), filter, options.Find().
		SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)+1))
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	posts, err := d.createPostList(cursor, askerID)

	if err != nil {
		log.Println(err)
		return nil, false, postErrors.DatabaseError(err.Error())
	}

	hasMore := len(posts) > limitConfig.Limit

	if hasMore {
		posts = posts[:len(posts)-1]
	}

	return posts, hasMore, err
}