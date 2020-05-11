package main

// Canvas Holds HTML5 canvas data
type Canvas struct {
	UUID string `json:"uuid" bson:"uuid"`
	Name string `json:"name" bson:"name"`
	Data string `json:"data" bson:"data"`
}
