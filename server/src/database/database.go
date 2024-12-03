package database

import (
	"context"

	postErrors "server/src/all_errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppDatabase struct {
	db *mongo.Database
}

func NewAppDatabase(client *mongo.Client) Database {
	return &AppDatabase{db: client.Database(DATABASE_NAME)}
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


