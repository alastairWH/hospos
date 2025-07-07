package db

import (
	"context"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance    *mongo.Client
	clientInstanceErr error
	mongoOnce         sync.Once
)

const dbName = "hospos"

// GetMongoClient returns a singleton MongoDB client
func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			uri = "mongodb://localhost:27017"
		}
		clientOptions := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			clientInstanceErr = err
			return
		}
		if err = client.Ping(context.Background(), nil); err != nil {
			clientInstanceErr = err
			return
		}
		clientInstance = client
		log.Println("Connected to MongoDB!")
	})
	return clientInstance, clientInstanceErr
}

// GetCollection returns a MongoDB collection by name
func GetCollection(name string) (*mongo.Collection, error) {
	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}
	return client.Database(dbName).Collection(name), nil
}
