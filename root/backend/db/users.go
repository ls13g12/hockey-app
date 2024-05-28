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

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plaintext password.
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
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


func CreateUser(db *mongo.Database, user User) error {
	coll := db.Collection("users")
	_, err := coll.InsertOne(context.TODO(), user)
	return err
}


func UpdateUser(db *mongo.Database, user User) error {
	coll := db.Collection("users")
	filter := bson.D{{Key: "user_id", Value: user.UserID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "username", Value: user.Username},
			{Key: "email", Value: user.Email},
			{Key: "password", Value: user.HashedPassword},
		}},
	}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	return err
}


func DeleteUser(db *mongo.Database, userID string) error {
	coll := db.Collection("users")
	filter := bson.D{{Key: "user_id", Value: userID}}
	_, err := coll.DeleteOne(context.TODO(), filter)
	return err
}
