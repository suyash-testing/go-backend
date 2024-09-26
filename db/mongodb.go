package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// ConnectMongo connects to the MongoDB database.
func ConnectMongo() *mongo.Client {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Could not connect to MongoDB:", err)
    }

    log.Println("Connected to MongoDB!")
    Client = client
    return client
}

// GetCollection gets the collection from the MongoDB database.
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    return client.Database("go-backend").Collection(collectionName)
}
