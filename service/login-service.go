package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type LoginService interface {
	Login(username string, password string) bool
}

type loginService struct {
	collection *mongo.Collection // MongoDB collection
}

func NewLoginService(client *mongo.Client, databaseName, collectionName string) LoginService {
	collection := client.Database(databaseName).Collection(collectionName)
	return &loginService{
		collection: collection,
	}
}

func (service *loginService) Login(email string, password string) bool {
	print("LoginService.Login() called\n")
	filter := bson.D{{Key: "email", Value: email}}

	var result struct {
		PasswordHash string `bson:"password"`
	}

	err := service.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		// Handle error (e.g., user not found)
		print("LoginService.Login() error: " + err.Error() + "\n")
		return false
	}

	// Compare stored password hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(password))
	if err != nil {
		// Passwords don't match
		print("LoginService.Login() passwords don't match\n")
		return false
	}

	// Passwords match
	return true
}
