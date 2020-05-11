package main

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/bson"
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

func getCanvases() ([]*Canvas, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.D{}, options.Find().SetProjection(bson.M{"name": 1, "uuid": 1}))

	if err != nil {
		log.Println("Failed reading collection ", err)
		return nil, err
	}
	defer cur.Close(ctx)
	res := new([]*Canvas)
	*res = make([]*Canvas, 0)
	cur.All(ctx, res)

	return *res, nil
}

func getCanvas(uuid string) (*Canvas, error) {
	canvas := &Canvas{}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if err := collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(canvas); err != nil {
		log.Println("failed decoding request ", err)
		return nil, err
	}

	return canvas, nil
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
