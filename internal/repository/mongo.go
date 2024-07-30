package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB() (*mongo.Client, *mongo.Collection, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongo://localhost:27017"))
	if err != nil {
		return nil, nil, err
	}
	db := client.Database("mafia")
	friends := db.Collection("friends")
	return client, friends, nil
}
