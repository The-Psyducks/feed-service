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

	frontPost := makeDBPostIntoFrontPost(newPost, false, false)

	return frontPost, err
}

func (d *AppDatabase) AddNewRetweet(newRetweet models.DBPost) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	_, err := postCollection.InsertOne(context.Background(), newRetweet)

	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	retweetCollection := d.db.Collection(RETWEET_COLLECTION)

	filter_original := bson.M{POST_ID_FIELD: newRetweet.Original_Post_ID}
	update := bson.M{"$inc": bson.M{RETWEET_FIELD: 1}}

	retweeter := bson.M{"$addToSet": bson.M{RETWEETERS_FIELD: newRetweet.Retweet_Author_ID}}

	_, err = postCollection.UpdateOne(context.Background(), filter_original, update)
	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	_, err = retweetCollection.UpdateOne(context.Background(), filter_original, retweeter, options.Update().SetUpsert(true))
	if err != nil {
		log.Println(err)
	}

	liked, err := d.hasLiked(newRetweet.Original_Post_ID, newRetweet.Retweet_Author_ID)

	post := makeDBPostIntoFrontPost(newRetweet, liked, true)

	return post, err
}

func (d *AppDatabase) GetPost(postID string, askerID string) (models.FrontPost, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	post, err := d.findPost(postID, postCollection)

	if err != nil {
		return models.FrontPost{}, err
	}

	liked, err := d.hasLiked(postID, askerID)

	retweeted, err_3 := d.hasRetweeted(post.Original_Post_ID, askerID)
	if err_3 != nil {
		return models.FrontPost{}, err_3
	}

	frontPost := makeDBPostIntoFrontPost(post, liked, retweeted)
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

func (d *AppDatabase) DeleteRetweet(postID string, userID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	retweetCollection := d.db.Collection(RETWEET_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$inc": bson.M{RETWEET_FIELD: -1}}
	retweeter := bson.M{"$pull": bson.M{RETWEETERS_FIELD: userID}}

	result, err := postCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return postErrors.ErrTwitsnapNotFound
	}

	_, err_2 := postCollection.UpdateOne(context.Background(), filter, update)

	if err_2 != nil {
		return err_2
	}

	_, err_3 := retweetCollection.UpdateOne(context.Background(), filter, retweeter)

	return err_3
}

func (d *AppDatabase) EditPost(postID string, editInfo models.EditPostExpectedFormat, askerID string) (models.FrontPost, error) {
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

	err_3 := d.updatePostPublic(postID, editInfo.Public)

	if err_3 != nil {
		return post, err_3
	}

	err_4 := d.updatePostMediaURL(postID, editInfo.MediaURL)

	if err_4 != nil {
		return post, err_4
	}

	dbPost, err = d.findPost(postID, postCollection)

	if err != nil {
		return post, err
	}

	retweeted, err_5 := d.hasRetweeted(postID, askerID)

	if err_5 != nil {
		return post, err_5
	}

	liked, err_5 := d.hasLiked(dbPost.Original_Post_ID, askerID)

	frontPost := makeDBPostIntoFrontPost(dbPost, liked, retweeted)

	return frontPost, err_5
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

	return err
}

func (d *AppDatabase) updatePostTags(postID string, newTags *[]string) error {

	if newTags == nil {
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

func (d *AppDatabase) updatePostMediaURL(postID string, newMediaURL *string) error {

	postCollection := d.db.Collection(FEED_COLLECTION)

	if newMediaURL == nil {
		return nil
	}

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{MEDIA_URL_FIELD: newMediaURL}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}

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

func (d *AppDatabase) GetUserFeedFollowing(following []string, askerID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)

	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}, "$or": []bson.M{
		{AUTHOR_ID_FIELD: bson.M{"$in": following}},
		{RETWEET_AUTHOR_FIELD: bson.M{"$in": following}},
	}}

	cursor, err := postCollection.Find(context.Background(), filter, options.Find().
		SetSort(bson.M{TIME_FIELD: -1}).SetSkip(int64(limitConfig.Skip)).SetLimit(int64(limitConfig.Limit)+1))
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())

	posts, err := d.createPostList(cursor, askerID)

	if err != nil {
		return nil, false, postErrors.DatabaseError(err.Error())
	}

	hasMore := len(posts) > limitConfig.Limit

	if hasMore {
		posts = posts[:len(posts)-1]
	}

	return posts, hasMore, err
}

func (d *AppDatabase) GetUserFeedInterests(interests []string, following []string, askerID string, limitConfig models.LimitConfig) ([]models.FrontPost, bool, error) {
	if len(interests) == 0 {
		return []models.FrontPost{}, false, postErrors.NoTagsFound()
	}

	postCollection := d.db.Collection(FEED_COLLECTION)
	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
	}

	filter := bson.M{TAGS_FIELD: bson.M{"$in": interests}, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}, "$or": []bson.M{
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

func (d *AppDatabase) GetUserFeedSingle(userId string, limitConfig models.LimitConfig, askerID string, following []string) ([]models.FrontPost, bool, error) {
	postCollection := d.db.Collection(FEED_COLLECTION)
	parsedTime, err := time.Parse(time.RFC3339, limitConfig.FromTime)

	if err != nil {
		log.Println(err)
		log.Println("User does not follow")
	}

	filter := bson.M{
		TIME_FIELD: bson.M{"$lt": parsedTime.UTC()},
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
			log.Println(word)
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

func (d *AppDatabase) LikeAPost(postID string, likerID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	likesCollection := d.db.Collection(LIKES_COLLECTION)

	liked, _ := d.hasLiked(postID, likerID)

	if liked {
		return postErrors.AlreadyLiked(postID)
	}

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$inc": bson.M{LIKES_FIELD: 1}}

	liker := bson.M{"$addToSet": bson.M{LIKERS_FIELD: likerID}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
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

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$inc": bson.M{LIKES_FIELD: -1}}

	liker := bson.M{"$pull": bson.M{LIKERS_FIELD: likerID}}

	_, err := postCollection.UpdateOne(context.Background(), filter, update)
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

func (d *AppDatabase) ClearDB() error {
	err := d.db.Collection(FEED_COLLECTION).Drop(context.Background())
	if err != nil {
		return postErrors.DatabaseError(err.Error())
	}
	err = d.db.Collection(LIKES_COLLECTION).Drop(context.Background())
	if err != nil {
		return postErrors.DatabaseError(err.Error())
	}

	err = d.db.Collection(RETWEET_COLLECTION).Drop(context.Background())
	if err != nil {
		return postErrors.DatabaseError(err.Error())
	}
	return nil
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
		liked, err_2 := d.hasLiked(dbPost.Post_ID, askerID)
		if err_2 != nil {
			return nil, err_2
		}
		retweeted, err_3 := d.hasRetweeted(dbPost.Original_Post_ID, askerID)
		if err_3 != nil {
			return nil, err_3
		}
		frontPost := makeDBPostIntoFrontPost(dbPost, liked, retweeted)
		posts = append(posts, frontPost)
	}

	return posts, err
}

func makeDBPostIntoFrontPost(post models.DBPost, liked bool, retweeted bool) models.FrontPost {
	author := models.AuthorInfo{
		Author_ID: post.Author_ID,
		Username:  "username",
		Alias:     "alias",
		PthotoURL: "photourl",
	}
	return models.NewFrontPost(post, author, liked, retweeted)
}

func (d *AppDatabase) hasLiked(postID string, likerID string) (bool, error) {
	if likerID == ADMIN {
		return false, nil
	}

	likesCollection := d.db.Collection(LIKES_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID, LIKERS_FIELD: likerID}

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

	filter := bson.M{POST_ID_FIELD: postID, RETWEETERS_FIELD: retweeterID}

	var res bson.M

	err := retweetCollection.FindOne(context.Background(), filter).Decode(&res)

	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	return err != mongo.ErrNoDocuments, nil
}
