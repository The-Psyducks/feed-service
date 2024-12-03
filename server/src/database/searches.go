package database

import (
	"context"
	"log"
	postErrors "server/src/all_errors"
	"server/src/models"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *AppDatabase) GetUserHashtags(interests []string, following []string, askerID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)

	if len(interests) == 0 {
		return []models.FrontPost{}, false, nil
	}

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{TAGS_FIELD: bson.M{"$all": interests}, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}, "$or": []bson.M{
		{PUBLIC_FIELD: true},
		{PUBLIC_FIELD: false, AUTHOR_ID_FIELD: bson.M{"$in": following}},
		{PUBLIC_FIELD: false, RETWEET_AUTHOR_FIELD: bson.M{"$in": following}},
	}}

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

func (d *AppDatabase) WordSearchPosts(words string, following []string, askerID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {

	postCollection := d.db.Collection(FEED_COLLECTION)

	filters := bson.A{}

	for _, word := range strings.Split(words, " ") {
		if word != "" {
			filters = append(filters, bson.M{CONTENT_FIELD: bson.M{"$regex": word, "$options": "i"}})
		}
	}

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{"$and": []bson.M{{"$or": filters}, {TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}}, {"$or": []bson.M{
		{PUBLIC_FIELD: true},
		{PUBLIC_FIELD: false, AUTHOR_ID_FIELD: bson.M{"$in": following}},
		{PUBLIC_FIELD: false, RETWEET_AUTHOR_FIELD: bson.M{"$in": following}},
	}}}}

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


func (d *AppDatabase) GetTrendingTopics() ([]string, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$" + TAGS_FIELD}}}},
		{{Key: "$match", Value: bson.D{
			{Key: TAGS_FIELD, Value: bson.D{{Key: "$type", Value: "string"}}},
			{Key: TIME_FIELD, Value: bson.D{{Key: "$type", Value: "date"}}},
		}}},
		{{Key: "$project",
			Value: bson.D{
				{Key: TAGS_FIELD, Value: 1},
				{Key: "timeDifference", Value: bson.D{
					{Key: "$divide", Value: bson.A{
						bson.D{{Key: "$subtract", Value: bson.A{
							bson.D{{Key: "$literal", Value: time.Now()}},
							"$" + TIME_FIELD,
						}}},
						1000 * 60 * 60,
					}},
				}},
			},
		}},
		{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: "$" + TAGS_FIELD},
				{Key: "totalOccurrences", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "averageTimeDifference", Value: bson.D{{Key: "$avg", Value: "$timeDifference"}}},
			},
		}},
		{{
			Key: "$project",
			Value: bson.D{
				{Key: TAGS_FIELD, Value: "$_id"},
				{Key: "score", Value: bson.D{
					{Key: "$multiply", Value: bson.A{
						"$totalOccurrences",
						bson.D{{Key: "$exp", Value: bson.D{
							{Key: "$multiply", Value: bson.A{-0.1, "$averageTimeDifference"}},
						}}},
					}},
				}},
			},
		}},
		{{Key: "$sort", Value: bson.D{{Key: "score", Value: -1}}}},
		{{Key: "$limit", Value: 20}},
	}

	cursor, err := postCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Println(err)
		return nil, postErrors.DatabaseError(err.Error())
	}

	var trendingTags []struct {
		Tag string `bson:"tags"`
	}

	if err = cursor.All(context.Background(), &trendingTags); err != nil {
		log.Println("Error decoding aggregation results:", err)
		return nil, postErrors.DatabaseError("Error decoding aggregation results")
	}

	tags := make([]string, len(trendingTags))
	for i, t := range trendingTags {
		tags[i] = t.Tag
	}

	return tags, nil
}