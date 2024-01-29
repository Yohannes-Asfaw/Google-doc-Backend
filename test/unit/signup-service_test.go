package unit_tests

import (
	"context"
	"testing"

	"github.com/khallihub/godoc/service"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignupServiceSuite struct {
	suite.Suite
	service service.SignupService
	client  *mongo.Client
}

func TestSignupServiceSuite(t *testing.T) {
	suite.Run(t, new(SignupServiceSuite))
}

func (s *SignupServiceSuite) SetupSuite() {
	// Setup MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Update with your MongoDB URI
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = client

	// Initialize the signup service
	s.service = service.NewSignupService(client, "testdb", "users")
}

func (s *SignupServiceSuite) SetupTest() {
	// Cleanup and prepare data before each test
	s.cleanupDatabase()
}

func (s *SignupServiceSuite) TearDownSuite() {
	// Close MongoDB connection after all tests
	if err := s.client.Disconnect(context.Background()); err != nil {
		s.T().Fatal(err)
	}
}

func (s *SignupServiceSuite) cleanupDatabase() {
	// Cleanup existing data in the test database
	_, err := s.client.Database("testdb").Collection("users").DeleteMany(context.Background(), bson.M{})
	if err != nil {
		s.T().Fatal(err)
	}
}
func (s *SignupServiceSuite) TestSignupValidUser() {
	username := "testuser"
	email := "testuser@test.com"
	password := "testpassword"

	// Call the method under test
	err := s.service.Signup(username, email, password)

	// Assertions
	s.NoError(err)

	// Verify that the user is inserted into the database
	filter := bson.D{{Key: "email", Value: email}}
	var insertedUser struct {
		Username string `bson:"username"`
		Email    string `bson:"email"`
	}

	err = s.client.Database("testdb").Collection("users").FindOne(context.Background(), filter).Decode(&insertedUser)
	s.NoError(err)
	s.Equal(username, insertedUser.Username)
	s.Equal(email, insertedUser.Email)
	// Add more assertions based on your use case
}



func (s *SignupServiceSuite) TestSignupDuplicateEmail() {
	// Prepare existing user
	existingUsername := "existinguser"
	existingEmail := "existinguser@test.com"
	existingPassword := "existingpassword"
	err := s.service.Signup(existingUsername, existingEmail, existingPassword)
	s.Require().NoError(err)

	// Attempt to register with the same email
	duplicateUsername := "duplicateuser"
	duplicateEmail := existingEmail
	duplicatePassword := "duplicatepassword"

	// Call the method under test
	err = s.service.Signup(duplicateUsername, duplicateEmail, duplicatePassword)

	// Assertions
	s.Error(err)
	s.EqualError(err, "username already exists")
}

// Add more test cases as needed
