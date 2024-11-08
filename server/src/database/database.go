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

	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	frontPost, err_2 := d.makeDBPostIntoFrontPost(newPost, newPost.Author_ID)

	return frontPost, err_2
}

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

	post, err := d.makeDBPostIntoFrontPost(newRetweet, newRetweet.Retweet_Author_ID)

	return post, err
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

func (d *AppDatabase) DeleteRetweet(postID string, userID string) error {
	postCollection := d.db.Collection(FEED_COLLECTION)
	retweetCollection := d.db.Collection(RETWEET_COLLECTION)

	filter := bson.M{ORIGINAL_POST_ID_FIELD: postID, RETWEET_AUTHOR_FIELD: userID}
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

	err_3 := d.updatePostPublic(postID, editInfo.Public)

	if err_3 != nil {
		return post, err_3
	}

	err_4 := d.updatePostMediaURL(postID, editInfo.MediaURL)

	if err_4 != nil {
		return post, err_4
	}

	dbPost, err_5 := d.findPost(postID, postCollection)

	if err_5 != nil {
		return post, err_5
	}

	frontPost, err_6 := d.makeDBPostIntoFrontPost(dbPost, askerID)

	return frontPost, err_6
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
	var mentions []string

	content :=  strings.Split(*newContent, " ")

	for _, word := range content {
		if strings.HasPrefix(word, "#") {
			tags = append(tags, word)
		}  else if strings.HasPrefix(word, "@") {
			mentions = append(mentions, word)
		}
	}

	err = d.updatePostTags(postID, &tags)

	if err != nil {
		log.Println(err)
	}

	err = d.updatePostMentions(postID, &mentions)

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

	fixedMentions := []string{}

	for _, word := range *newMentions {
		fixedMentions = append(fixedMentions, word[1:])
	}

	postCollection := d.db.Collection(FEED_COLLECTION)

	filter := bson.M{POST_ID_FIELD: postID}
	update := bson.M{"$set": bson.M{MENTIONS_FIELD: fixedMentions}}

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

	log.Println(interests)

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

	log.Println(posts)

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
				{AUTHOR_ID_FIELD: userId, IS_RETWEET_FIELD: false},
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
