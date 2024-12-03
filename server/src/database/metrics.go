package database

import (
	"context"
	"log"
	"server/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *AppDatabase) GetUserMetrics(userID string, limits models.MetricLimits) (models.UserMetrics, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)

	parsedFromTime, err := time.Parse(time.RFC3339, limits.FromTime)
	if err != nil {
		log.Println(err)
		return models.UserMetrics{}, err
	}

	parsedToTime, err := time.Parse(time.RFC3339, limits.ToTime)
	if err != nil {
		log.Println(err)
		return models.UserMetrics{}, err
	}

	pipeline := mongo.Pipeline{

		bson.D{{Key: "$match", Value: bson.D{
			{Key: TIME_FIELD, Value: bson.D{{Key: "$gte", Value: parsedFromTime.UTC()}, {Key: "$lt", Value: parsedToTime.UTC()}}},
			{Key: AUTHOR_ID_FIELD, Value: userID},
			{Key: IS_RETWEET_FIELD, Value: false},
		}}},

		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "likes", Value: bson.D{{Key: "$sum", Value: "$" + LIKES_FIELD}}},
			{Key: "retweets", Value: bson.D{{Key: "$sum", Value: "$" + RETWEET_FIELD}}},
			{Key: "posts", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := postCollection.Aggregate(context.Background(), pipeline)

	if err != nil {
		log.Println(err)
		return models.UserMetrics{}, err
	}

	var result []bson.M

	if err := cursor.All(context.Background(), &result); err != nil {
		log.Fatal(err)
	}

	var metrics models.UserMetrics

	if len(result) > 0 {

		metrics = models.UserMetrics{Likes: convertToInt(result[0]["likes"]), Retweets: convertToInt(result[0]["retweets"]), Posts: convertToInt(result[0]["posts"])}
	} else {
		metrics = models.UserMetrics{Likes: 0, Retweets: 0, Posts: 0}
	}

	return metrics, nil
}