package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB(dsn string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))

	if err != nil {
		return nil, err
	}

	return client, client.Ping(context.TODO(), nil)
}
