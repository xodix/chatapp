package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func configDB() *mongo.Database {
	mongoURI := os.Getenv("MONGODB_URI")

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB:", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln("Failed to ping MongoDB:", err)
	}
	db := client.Database("chatapp")
	if db == nil {
		log.Fatalln("Failed to access MongoDB Database")
	}

	return db
}
