package main

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectedToMongo = make(chan struct{})
var collection *mongo.Collection

func startDb() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Println(err)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), time.Second)
	for client.Ping(ctx, readpref.Primary()) != nil {
		time.Sleep(time.Second)
	}

	collection = client.Database(viper.GetString("mongoDbName")).Collection(viper.GetString("mongoCollectionName"))
	close(connectedToMongo)
}

func saveCanvas(canvas *Canvas) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// buf, err := bson.Marshal(canvas)
	// if err != nil {
	// 	log.Println("failed marshaling canvas ", err)
	// 	return err
	// }

	if _, err := collection.InsertOne(ctx, canvas); err != nil {
		log.Println("could not persist canvas ", err)
		return err
	}

	return nil
}
