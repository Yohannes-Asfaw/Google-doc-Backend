package service

import (
	"context"
	"fmt"
	// "errors"
	// "fmt"
	"github.com/khallihub/godoc/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocumentService interface {
	GetAllDocuments(email dto.Email) ([]*dto.Document, error)
	SearchDocuments(email string, searchQuery string) ([]*dto.Document, error)
	CreateDocument(author string, title string, body interface{}, readAccess []string, writeAccess []string) (string, error)
	UpdateDocument(documentID string, body dto.DocumentData) error
	GetDocumentByID(documentID string) (*dto.Document, error)
	UpdateTitle(documentID string, title string) (string, error)
	UpdateCollaborators(documentID string, collaborators dto.Access) (dto.Document, error)
	DeleteDocument(documentID string) (bool, error)
}

type documentService struct {
	collection *mongo.Collection // MongoDB collection
}

func NewDocumentService(client *mongo.Client, databaseName, collectionName string) DocumentService {
	collection := client.Database(databaseName).Collection(collectionName)
	return &documentService{
		collection: collection,
	}
}

func (service *documentService) GetAllDocuments(email dto.Email) ([]*dto.Document, error) {
	filter := bson.M{"$or": []bson.M{{"author": email.Email}, {"readAccess": email.Email}, {"writeAccess": email.Email}}}
	cursor, err := service.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var documents []*dto.Document
	for cursor.Next(context.Background()) {
		var document dto.Document
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}
		documentDTO := &dto.Document{
			ID:    document.ID,
			Title: document.Title,
		}

		documents = append(documents, documentDTO)
	}
	return documents, nil
}

func (service *documentService) SearchDocuments(email string, searchQuery string) ([]*dto.Document, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"author": email},
			{"readAccess": email},
			{"writeAccess": email},
		},
		"title": bson.M{
			"$regex":   searchQuery,
			"$options": "i",
		},
	}
	cursor, err := service.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var documents []*dto.Document
	for cursor.Next(context.Background()) {
		var document dto.Document
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}
		documentDTO := &dto.Document{
			ID:    document.ID,
			Title: document.Title,
		}

		documents = append(documents, documentDTO)
	}
	return documents, nil
}

func (service *documentService) CreateDocument(author string, title string, body interface{}, readAccess []string, writeAccess []string) (string, error) {
	// Implement logic to create a document in the MongoDB collection
	newDocument := bson.D{
		{Key: "author", Value: author},
		{Key: "readAccess", Value: readAccess},
		{Key: "writeAccess", Value: writeAccess},
		{Key: "title", Value: title},
		{Key: "body", Value: body},
	}

	document, err := service.collection.InsertOne(context.Background(), newDocument)
	if err != nil {
		return "", err
	}

	insertedID, ok := document.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert InsertedID to string")
	}

	return insertedID.Hex(), nil
}

func (service *documentService) UpdateDocument(documentID string, incomingData dto.DocumentData) error {
	objectID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"data.ops": incomingData.Ops}}
	filter := bson.M{"_id": objectID}
	_, err = service.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (service *documentService) GetDocumentByID(documentID string) (*dto.Document, error) {
	// Implement logic to fetch a document by its ID from the MongoDB collection
	var document dto.Document
	objectID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	err = service.collection.FindOne(context.Background(), filter).Decode(&document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (service *documentService) UpdateTitle(documentID string, title string) (string, error) {
	// Implement logic to update the title of a document in the MongoDB collection
	objectID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return "", err
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"title": title}}
	_, err = service.collection.UpdateOne(context.Background(), filter, update)
	return "", err
}

func (service *documentService) UpdateCollaborators(documentID string, collaborators dto.Access) (dto.Document, error) {
    objectID, err := primitive.ObjectIDFromHex(documentID)
    if err != nil {
        fmt.Print(err)
        return dto.Document{}, err
    }

    filter := bson.M{"_id": objectID}
    update := bson.M{"$set": bson.M{"readAccess": collaborators.ReadAccess, "writeAccess": collaborators.WriteAccess}}

    // Perform the update operation
    _, err = service.collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return dto.Document{}, err
    }

    // Retrieve the updated document
    var updatedDocument dto.Document
    err = service.collection.FindOne(context.Background(), filter).Decode(&updatedDocument)
    if err != nil {
        return dto.Document{}, err
    }
    return updatedDocument, nil
}

func (service *documentService) DeleteDocument(documentID string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": objectID}
	result, err := service.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}
