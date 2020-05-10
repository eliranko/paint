package main

// Canvas Holds HTML5 canvas data
type Canvas struct {
	UUID string `bson:"uuid"`
	Data string `bson:"data"`
}
