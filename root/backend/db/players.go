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

func AllPlayers(db *mongo.Database) ([]Player, error){
	var players []Player

	coll := db.Collection("players") //Set db name to command line flag
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

func GetPlayer(db *mongo.Database, playerID string) (Player, error) {
	var player Player

	coll := db.Collection("players")
	filter := bson.D{{Key: "player_id", Value: playerID}}
	err := coll.FindOne(context.TODO(), filter).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return player, nil // Return empty player if not found
		}
		return player, err
	}

	return player, nil
}

func CreatePlayer(db *mongo.Database, player Player) error {
	coll := db.Collection("players")
	player.Created = time.Now()
	_, err := coll.InsertOne(context.TODO(), player)
	return err
}

func UpdatePlayer(db *mongo.Database, player Player) error {
	coll := db.Collection("players")
	filter := bson.D{{Key: "player_id", Value: player.PlayerID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "first_name", Value: player.FirstName},
			{Key: "last_name", Value: player.LastName},
			{Key: "nickname", Value: player.Nickname},
			{Key: "home_shirt_number", Value: player.HomeShirtNumber},
			{Key: "date_of_birth", Value: player.DateOfBirth},
		}},
	}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}

func DeletePlayer(db *mongo.Database, playerID string) error {
	coll := db.Collection("players")
	filter := bson.D{{Key: "player_id", Value: playerID}}
	_, err := coll.DeleteOne(context.TODO(), filter)
	return err
}
