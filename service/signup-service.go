package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignupService interface {
	Signup(username string, email string, password string) error
}

type signupService struct {
	collection *mongo.Collection // MongoDB collection
}

func NewSignupService(client *mongo.Client, databaseName, collectionName string) SignupService {
	collection := client.Database(databaseName).Collection(collectionName)
	return &signupService{
		collection: collection,
	}
}

func (service *signupService) Signup(username string, email string, password string) error {
	print("SignupService.Signup called\n") 

	// Check if the username already exists
	filter := bson.D{{Key: "email", Value: email}}
	var existingUser struct{}
	err := service.collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		return errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create a new user document
	newUser := bson.D{
		{Key: "username", Value: username},
		{Key: "email", Value: email},
		{Key: "password", Value: string(hashedPassword)},
	}

	// Insert the new user into the database
	_, err = service.collection.InsertOne(context.Background(), newUser)
	if err != nil {
		return err
	}

	return nil
}
