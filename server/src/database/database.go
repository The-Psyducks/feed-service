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

	newPostRetweet, err := d.findPost(newRetweet.Post_ID, postCollection)

	if err != nil {
		log.Println(err)
		return models.FrontPost{}, err
	}

	post, err := d.makeDBPostIntoFrontPost(newPostRetweet, newPostRetweet.Retweet_Author_ID)

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

	filter := bson.M{TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}, BLOCKED_FIELD: false, "$or": []bson.M{
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

	filter := bson.M{TAGS_FIELD: bson.M{"$in": interests}, TIME_FIELD: bson.M{"$lt": parsedTime.UTC()}, BLOCKED_FIELD: false, "$or": []bson.M{
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
		TIME_FIELD:    bson.M{"$lt": parsedTime.UTC()},
		BLOCKED_FIELD: false,
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
