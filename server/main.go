package main

import (
	"fmt"
	"log"
	"time"
	"context"
	"server/src/database"
	"server/src/router"

	"server/src"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	config := src.ConfigEnv()

	gin.SetMode(config.Gin_Mode)

	gin.ForceConsoleColor()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo_URI))

	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := database.NewDatabase(client)

	r := router.CreateRouter(db)

	address := fmt.Sprintf("%s:%s", config.Host, config.Port)

	if err := r.Run(address); err != nil {
		log.Fatal("Error running server: ", err)
	}
}
