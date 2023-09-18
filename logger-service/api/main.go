package main

import (
	"context"
	"github.com/Noah-Wilderom/cloud-services/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
	"time"
)

const (
	grpcPort     = 5001
	mongoURL     = "mongodb://mongo:27017"
	mongoTimeout = 15 * time.Second
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongodb
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Create a WaitGroup
	var wg sync.WaitGroup

	// Increment the WaitGroup to indicate a goroutine is starting
	wg.Add(1)

	// Start your gRPC server in a goroutine
	go func() {
		defer wg.Done() // Decrement the WaitGroup when done
		app.gRPCListen()
	}()

	// Wait for the gRPC server to start
	wg.Wait()
}

func connectToMongo() (*mongo.Client, error) {
	// create the connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})

	// connect
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connection:", err)
		return nil, err
	}

	log.Println("Connected to mongodb")

	return conn, nil
}
