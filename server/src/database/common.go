package database

import (
	"context"
	"log"
	postErrors "server/src/all_errors"
	"server/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *AppDatabase) findPost(postID string, postCollection *mongo.Collection) (models.DBPost, error) {
	var post models.DBPost
	filter := bson.M{POST_ID_FIELD: postID, BLOCKED_FIELD: false}
	err := postCollection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		log.Println(err)
		err = postErrors.ErrTwitsnapNotFound
	}
	return post, err
}

func (d *AppDatabase) createPostList(cursor *mongo.Cursor, askerID string) ([]models.FrontPost, error) {
	var posts []models.FrontPost
	var err error

	for cursor.Next(context.Background()) {
		var dbPost models.DBPost
		err = cursor.Decode(&dbPost)
		if err != nil {
			return nil, err
		}

		frontPost, err_2 := d.makeDBPostIntoFrontPost(dbPost, askerID)

		if err_2 != nil {
			return nil, err_2
		}

		posts = append(posts, frontPost)
	}

	return posts, err
}

func (d *AppDatabase) makeDBPostIntoFrontPost(post models.DBPost, askerID string) (models.FrontPost, error) {
	author := models.AuthorInfo{
		Author_ID: post.Author_ID,
		Username:  "username",
		Alias:     "alias",
		PthotoURL: "photourl",
	}

	if len(post.Tags) == 0 {
		post.Tags = []string{}
	}

	if len(post.Mentions) == 0 {
		post.Mentions = []string{}
	}

	liked, err_2 := d.hasLiked(post.Original_Post_ID, askerID)
	if err_2 != nil {
		return models.FrontPost{}, err_2
	}

	retweeted, err_3 := d.hasRetweeted(post.Original_Post_ID, askerID)
	if err_3 != nil {
		return models.FrontPost{}, err_3
	}

	bookmarked, err_4 := d.hasBookmark(post.Post_ID, askerID)
	if err_4 != nil {
		return models.FrontPost{}, err_4
	}
	return models.NewFrontPost(post, author, liked, retweeted, bookmarked), nil
}

func (d *AppDatabase) hasLiked(postID string, likerID string) (bool, error) {
	if likerID == ADMIN {
		return false, nil
	}

	likesCollection := d.db.Collection(LIKES_COLLECTION)

	filter := bson.M{ORIGINAL_POST_ID_FIELD: postID, LIKERS_FIELD: likerID}

	var res bson.M

	err := likesCollection.FindOne(context.Background(), filter).Decode(&res)

	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	return err != mongo.ErrNoDocuments, nil
}

func (d *AppDatabase) hasRetweeted(postID string, retweeterID string) (bool, error) {
	if retweeterID == ADMIN {
		return false, nil
	}

	retweetCollection := d.db.Collection(RETWEET_COLLECTION)

	filter := bson.M{ORIGINAL_POST_ID_FIELD: postID, RETWEETERS_FIELD: retweeterID}

	var res bson.M

	err := retweetCollection.FindOne(context.Background(), filter).Decode(&res)

	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	return err != mongo.ErrNoDocuments, nil
}

func (d *AppDatabase) hasBookmark(postID string, userID string) (bool, error) {
	favoritesCollection := d.db.Collection(BOOKMARK_COLLECTION)

	filter := bson.M{AUTHOR_ID_FIELD: userID, POST_ID_FIELD: postID}

	var res bson.M

	err := favoritesCollection.FindOne(context.Background(), filter).Decode(&res)

	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	return err != mongo.ErrNoDocuments, nil
}

func convertToInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	default:
		return 0
	}
}