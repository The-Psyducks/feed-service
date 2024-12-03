package database

import (
	"context"
	"log"
	"server/src/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)


func (d *AppDatabase) EditPost(postID string, editInfo models.EditPostExpectedFormat, askerID string) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	var post models.FrontPost
	var dbPost models.DBPost

	err := d.updatePostContent(postID, editInfo.Content)

	if err != nil {
		return post, err
	}

	err_3 := d.updatePostPublic(postID, editInfo.Public)

	if err_3 != nil {
		return post, err_3
	}

	err_4 := d.updatePostMediaURL(postID, editInfo.MediaInfo)

	if err_4 != nil {
		return post, err_4
	}

	err_5 := d.updatePostMentions(postID, &editInfo.Mentions)

	if err_5 != nil {
		return post, err_5
	}

	dbPost, err_6 := d.findPost(postID, postCollection)

	if err_6 != nil {
		return post, err_6
	}

	frontPost, err_7 := d.makeDBPostIntoFrontPost(dbPost, askerID)

	return frontPost, err_7
}

func (d *AppDatabase) updatePostContent(postID string, newContent *string) error {

	if newContent == nil {
		return nil
	}

	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{CONTENT_FIELD: newContent}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	var tags []string

	content := strings.Split(*newContent, " ")

	for _, word := range content {
		if strings.HasPrefix(word, "#") {
			tags = append(tags, word)
		}
	}

	err = d.updatePostTags(postID, &tags)

	return err
}

func (d *AppDatabase) updatePostTags(postID string, newTags *[]string) error {

	if newTags == nil {
		return nil
	}

	fixedTags := []string{}

	for _, word := range *newTags {
		fixedTags = append(fixedTags, word[1:])
	}

	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{TAGS_FIELD: fixedTags}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) updatePostMentions(postID string, newMentions *[]string) error {

	if newMentions == nil {
		return nil
	}

	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{MENTIONS_FIELD: newMentions}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) updatePostPublic(postID string, newPublic *bool) error {

	postCollection := d.db.Collection(FEED_COLLECTION)

	if newPublic == nil {
		return nil
	}

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{PUBLIC_FIELD: newPublic}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (d *AppDatabase) updatePostMediaURL(postID string, newMediaInfo *models.MediaInfo) error {

	postCollection := d.db.Collection(FEED_COLLECTION)

	if newMediaInfo == nil {
		return nil
	}

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{MEDIA_INFO_FIELD: newMediaInfo}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

	return err
}