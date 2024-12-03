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

func (d *AppDatabase) AddFavorite(postID string, userID string) error {
	favoritesCollection := d.db.Collection(BOOKMARK_COLLECTION)

	filter := bson.M{AUTHOR_ID_FIELD: userID}
	update := bson.M{"$addToSet": bson.M{POST_ID_FIELD: postID}}

	_, err := favoritesCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))

	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) RemoveFavorite(postID string, userID string) error {
	favoritesCollection := d.db.Collection(BOOKMARK_COLLECTION)

	filter := bson.M{AUTHOR_ID_FIELD: userID}
	update := bson.M{"$pull": bson.M{POST_ID_FIELD: postID}}

	_, err := favoritesCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) GetUserFavorites(userID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {

	favoritesCollection := d.db.Collection(BOOKMARK_COLLECTION)
	postCollection := d.db.Collection(FEED_COLLECTION)

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{AUTHOR_ID_FIELD: userID}

	cursor, err := favoritesCollection.Find(context.Background(), filter, options.Find())

	if err != nil {
		log.Println(err)
		return nil, false, postErrors.DatabaseError(err.Error())
	}

	postIDs := []string{}

	for cursor.Next(context.Background()) {
		var res bson.M
		err = cursor.Decode(&res)
		if err != nil {
			log.Println(err)
			return nil, false, postErrors.DatabaseError(err.Error())
		}

		if postIDArray, ok := res[POST_ID_FIELD].(bson.A); ok {
			for _, postID := range postIDArray {
				if idStr, ok := postID.(string); ok {
					postIDs = append(postIDs, idStr)
				} else {
					log.Printf("Unexpected post ID type: %T", postID)
				}
			}
		} else {
			log.Printf("Unexpected type for POST_ID_FIELD: %T", res[POST_ID_FIELD])
		}
	}

	filter = bson.M{POST_ID_FIELD: bson.M{"$in": postIDs}, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}}

	cursor, err = postCollection.Find(context.Background(), filter, options.Find().
		SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)+1))

	if err != nil {
		log.Println(err)
		return nil, false, postErrors.DatabaseError(err.Error())
	}

	posts, err := d.createPostList(cursor, userID)

	if err != nil {
		log.Println(err)
		return nil, false, postErrors.DatabaseError(err.Error())
	}

	hasMore := len(posts) > limitConfig.Limit

	return posts, hasMore, nil
}