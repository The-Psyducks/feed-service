package database

import (
	"context"
	"log"


	"go.mongodb.org/mongo-driver/bson"
)

func (d *AppDatabase) BlockPost(postId string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postId}
	update := bson.M{"$set": bson.M{BLOCKED_FIELD: true}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) UnBlockPost(postId string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postId}
	update := bson.M{"$set": bson.M{BLOCKED_FIELD: false}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Println(err)
	}

	return err
}