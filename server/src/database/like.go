package database

import (
	"context"
	"log"
	postErrors "server/src/all_errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *AppDatabase) LikeAPost(postID string, likerID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	likesCollection := d.db.Collection(LIKES_COLLECTION)

	liked, _ := d.hasLiked(postID, likerID)

	if liked {
		return postErrors.AlreadyLiked(postID)
	}

	filter := bson.M{ORIGINAL_POST_ID_FIELD: postID}
	update := bson.M{"$inc": bson.M{LIKES_FIELD: 1}}

	liker := bson.M{"$addToSet": bson.M{LIKERS_FIELD: likerID}}

	_, err := postCollection.UpdateMany(context.Background(), filter, update)
	
	if err != nil {
		log.Println(err)
		return postErrors.TwitsnapNotFound(postID)
	}

	_, err = likesCollection.UpdateOne(context.Background(), filter, liker, options.Update().SetUpsert(true))
	if err != nil {
		log.Println(err)
		return postErrors.TwitsnapNotFound(postID)
	}

	return nil
}

func (d *AppDatabase) UnLikeAPost(postID string, likerID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	likesCollection := d.db.Collection(LIKES_COLLECTION)

	filter := bson.M{ORIGINAL_POST_ID_FIELD: postID}
	update := bson.M{"$inc": bson.M{LIKES_FIELD: -1}}

	liker := bson.M{"$pull": bson.M{LIKERS_FIELD: likerID}}

	_, err := postCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = likesCollection.UpdateOne(context.Background(), filter, liker)

	if err != nil {
		log.Println(err)
	}

	return err
}
