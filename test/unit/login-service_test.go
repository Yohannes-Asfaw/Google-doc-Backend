package unit_tests

import (
	"context"
	"testing"

	"github.com/khallihub/godoc/service"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type LoginServiceSuite struct {
	suite.Suite
	service service.LoginService
	client  *mongo.Client
}

func TestLoginServiceSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceSuite))
}

func (s *LoginServiceSuite) SetupSuite() {
	// Setup MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Update with your MongoDB URI
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = client

	// Initialize the login service
	s.service = service.NewLoginService(client, "testdb", "users")
}

func (s *LoginServiceSuite) SetupTest() {
	// Cleanup and prepare data before each test
	s.cleanupDatabase()
	s.prepareTestData()
}

func (s *LoginServiceSuite) TearDownSuite() {
	// Close MongoDB connection after all tests
	if err := s.client.Disconnect(context.Background()); err != nil {
		s.T().Fatal(err)
	}
}

func (s *LoginServiceSuite) cleanupDatabase() {
	// Cleanup existing data in the test database
	_, err := s.client.Database("testdb").Collection("users").DeleteMany(context.Background(), bson.M{})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *LoginServiceSuite) prepareTestData() {
	// Insert test data into the database for each test case
	// You can insert users relevant to your test cases
	// Example:
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	user := bson.D{{Key: "email", Value: "testuser@test.com"}, {Key: "password", Value: hashedPassword}}

	_, err := s.client.Database("testdb").Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *LoginServiceSuite) TestLoginValidCredentials() {
	username := "testuser@test.com"
	password := "testpassword"

	// Call the method under test
	result := s.service.Login(username, password)

	// Assertions
	s.True(result)
}

func (s *LoginServiceSuite) TestLoginInvalidUsername() {
	username := "invaliduser@test.com"
	password := "testpassword"

	// Call the method under test
	result := s.service.Login(username, password)

	// Assertions
	s.False(result)
}

func (s *LoginServiceSuite) TestLoginInvalidPassword() {
	username := "testuser@test.com"
	password := "invalidpassword"

	// Call the method under test
	result := s.service.Login(username, password)

	// Assertions
	s.False(result)
}
