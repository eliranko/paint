package main

import (
	"context"
	"errors"
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
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongoUrl")))
	if err != nil {
		log.Panic(err)
	}

	timeout := 30 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err = client.Connect(ctx); err != nil {
		log.Println(err)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), timeout)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println(err)
		return
	}

	log.Println("Connected to db at ", viper.GetString("mongoUrl"))
	collection = client.Database(viper.GetString("mongoDbName")).Collection(viper.GetString("mongoCollectionName"))
	close(connectedToMongo)
}

func getCanvases(ctx context.Context) ([]*Canvas, error) {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return nil, errors.New("not connected to db")
	}

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

func getCanvas(ctx context.Context, uuid string) (*Canvas, error) {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return nil, errors.New("not connected to db")
	}

	canvas := &Canvas{}
	if err := collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(canvas); err != nil {
		log.Println("failed decoding request ", err)
		return nil, err
	}

	return canvas, nil
}

func saveCanvas(ctx context.Context, canvas *Canvas) error {
	select {
	case <-connectedToMongo:
	case <-ctx.Done():
		return errors.New("not connected to db")
	}

	if _, err := collection.InsertOne(ctx, canvas); err != nil {
		log.Println("could not persist canvas ", err)
		return err
	}

	return nil
}
