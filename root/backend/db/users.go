package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	ERR_EMAIL_ALREADY_EXISTS = "Email already exists"
	ERR_AUTH_USER_NOT_FOUND = "Authenticated Failed - User not found"
)

type User struct {
	UserID    string    `bson:"user_id" json:"user_id"`
	Username  string    `bson:"username" json:"username,required"`
	Email     string    `bson:"email" json:"email,required"`
	HashedPassword  string    `bson:"hashed_password" json:"hashed_password"`
	Created   time.Time `bson:"created" json:"created"`
}


func GetUser(db *mongo.Database, userID string) (User, error) {
	var user User

	coll := db.Collection("users")
	filter := bson.D{{Key: "user_id", Value: userID}}
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func Exists(db *mongo.Database, email string) (bool, error) {
	coll := db.Collection("users")
	filter := bson.D{{Key: "email", Value: email}}

	var result bson.M
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func Authenticate(db *mongo.Database, email string, password string) (string, error) {
	coll := db.Collection("users")
	filter := bson.D{{Key: "email", Value: email}}

	var foundUser User
	err := coll.FindOne(context.TODO(), filter).Decode(&foundUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.HashedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	return foundUser.UserID, nil
}


func CreateUser(db *mongo.Database, user User) error {
	coll := db.Collection("users")
	_, err := coll.InsertOne(context.TODO(), user)
	return err
}
