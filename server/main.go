package main

import (
	"fmt"
	"log"
	"time"
	"os"
	"context"
	"server/src/database"
	"server/src/router"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	gin.SetMode(os.Getenv("GIN_MODE"))

	gin.ForceConsoleColor()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := database.NewAppDatabase(client)

	// err = db.ClearDB()

	// if err != nil {
	// 	log.Fatal("Error clearing database: ", err)
	// }

	r := router.CreateRouter(db)

	address := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	log.Println("Server running on: ", address)

	if err := r.Run(address); err != nil {
		log.Fatal("Error running server: ", err)
	}
}
