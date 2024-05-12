package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Player struct {
	PlayerID         string    `bson:"player_id" json:"player_id"`
	FirstName        string    `bson:"first_name" json:"first_name"`
	LastName         string    `bson:"last_name" json:"last_name"`
	Nickname         string    `bson:"nickname" json:"nickname"`
	HomeShirtNumber  int       `bson:"home_shirt_number" json:"home_shirt_number"`
	DateOfBirth      time.Time `bson:"date_of_birth" json:"date_of_birth"`
	Created          time.Time `bson:"created" json:"created"`
}

func AllPlayers(dbClient *mongo.Client) ([]Player, error){
	var players []Player

	coll := dbClient.Database("hockeydb").Collection("players") //Set db name to command line flag
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}
