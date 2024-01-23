package unit_tests

import (
	"context"
	"testing"

	"github.com/khallihub/godoc/dto"
	"github.com/khallihub/godoc/service"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentServiceSuite struct {
	suite.Suite
	service service.DocumentService
	client  *mongo.Client
}

func TestDocumentServiceSuite(t *testing.T) {
	suite.Run(t, new(DocumentServiceSuite))
}

func (s *DocumentServiceSuite) SetupSuite() {
	// Setup MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Update with your MongoDB URI
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = client

	// Initialize the document service
	s.service = service.NewDocumentService(client, "testdb", "documentCollection")
}

func (s *DocumentServiceSuite) SetupTest() {
	// Cleanup and prepare data before each test
	s.cleanupDatabase()
	s.prepareTestData()
}

func (s *DocumentServiceSuite) TearDownSuite() {
	// Close MongoDB connection after all tests
	if err := s.client.Disconnect(context.Background()); err != nil {
		s.T().Fatal(err)
	}
}

func (s *DocumentServiceSuite) cleanupDatabase() {
	// Cleanup existing data in the test database
	_, err := s.client.Database("testdb").Collection("documentCollection").DeleteMany(context.Background(), bson.M{})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *DocumentServiceSuite) prepareTestData() {
	_, err := s.service.CreateDocument("author@test.com", "Test Document", "Test body", []string{"read@test.com"}, []string{"write@test.com"})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *DocumentServiceSuite) TestGetAllDocuments() {
	email := dto.Email{Email: "author@test.com"}

	// Call the method under test
	documents, err := s.service.GetAllDocuments(email)

	// Assertions
	s.NoError(err)
	s.Len(documents, 1)
	s.Equal("Test Document", documents[0].Title)
}

func (s *DocumentServiceSuite) TestSearchDocuments() {
	email := "author@test.com"
	searchQuery := "Test"

	// Call the method under test
	documents, err := s.service.SearchDocuments(email, searchQuery)

	// Assertions
	s.NoError(err)
	s.Len(documents, 1)
	s.Equal("Test Document", documents[0].Title)
	// Add more assertions based on your use case
}

func (s *DocumentServiceSuite) TestCreateDocument() {
	author := "new_author@test.com"
	title := "New Test Document"
	body := "New Test Body"
	readAccess := []string{"read@test.com"}
	writeAccess := []string{"write@test.com"}

	// Call the method under test
	documentID, err := s.service.CreateDocument(author, title, body, readAccess, writeAccess)

	// Assertions
	s.NoError(err)
	s.NotEmpty(documentID)
	// Add more assertions based on your use case
}

func (s *DocumentServiceSuite) TestUpdateDocument() {
	// Prepare test data
	documents, err := s.service.GetAllDocuments(dto.Email{Email: "author@test.com"})
	s.Require().NoError(err)
	s.Require().Len(documents, 1)
	documentID := documents[0].ID

	// Prepare update data
	incomingData := dto.DocumentData{
		Ops: []map[string]interface{}{
			{
				"insert": map[string]interface{}{
					"position": 0,
					"text":     "Updated Text",
				},
			},
		},
	}

	// Call the method under test
	err = s.service.UpdateDocument(documentID, incomingData)

	// Assertions
	s.NoError(err)
	// Add more assertions based on your use case
}

func (s *DocumentServiceSuite) TestGetDocumentByID() {
	// Prepare test data
	documents, err := s.service.GetAllDocuments(dto.Email{Email: "author@test.com"})
	s.Require().NoError(err)
	s.Require().Len(documents, 1)
	documentID := documents[0].ID

	// Call the method under test
	document, err := s.service.GetDocumentByID(documentID)

	// Assertions
	s.NoError(err)
	s.NotNil(document)
	s.Equal(documentID, document.ID)
	// Add more assertions based on your use case
}

func (s *DocumentServiceSuite) TestUpdateTitle() {
	// Prepare test data
	documents, err := s.service.GetAllDocuments(dto.Email{Email: "author@test.com"})
	s.Require().NoError(err)
	s.Require().Len(documents, 1)
	documentID := documents[0].ID

	// Prepare update title
	newTitle := "Updated Title"

	// Call the method under test
	updatedTitle, err := s.service.UpdateTitle(documentID, newTitle)

	// Assertions
	s.NoError(err)
	s.Equal("", updatedTitle)
	// Add more assertions based on your use case
}

func (s *DocumentServiceSuite) TestUpdateCollaborators() {
	// Prepare test data
	documents, err := s.service.GetAllDocuments(dto.Email{Email: "author@test.com"})
	s.Require().NoError(err)
	s.Require().Len(documents, 1)
	documentID := documents[0].ID

	// Prepare collaborators
	newCollaborators := dto.Access{
		ReadAccess:  []string{"new_read@test.com"},
		WriteAccess: []string{"new_write@test.com"},
	}

	// Call the method under test
	updatedDocument, err := s.service.UpdateCollaborators(documentID, newCollaborators)

	// Assertions
	s.NoError(err)
	s.NotNil(updatedDocument)
	s.Equal(newCollaborators.ReadAccess, updatedDocument.ReadAccess)
	s.Equal(newCollaborators.WriteAccess, updatedDocument.WriteAccess)
	// Add more assertions based on your use case
}

func (s *DocumentServiceSuite) TestDeleteDocument() {
	// Prepare test data
	documents, err := s.service.GetAllDocuments(dto.Email{Email: "author@test.com"})
	s.Require().NoError(err)
	s.Require().Len(documents, 1)
	documentID := documents[0].ID

	// Declare the variable "deleted"
	var deleted bool

	// Call the method under test
	deleted, err = s.service.DeleteDocument(documentID)

	// Assertions
	s.NoError(err)
	s.Assert().Equal(true, deleted)

	// Verify that the document is deleted
	deletedDocument, err := s.service.GetDocumentByID(documentID)
	s.Error(err)
	s.Nil(deletedDocument)
}