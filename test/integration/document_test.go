package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/khallihub/godoc/service"
	"github.com/stretchr/testify/suite"
)

type DocumentEndpointsSuite struct {
	suite.Suite
	jwtService service.JWTService
}

func TestDocumentEndpointsSuite(t *testing.T) {
	suite.Run(t, new(DocumentEndpointsSuite))
}

func (s *DocumentEndpointsSuite) SetupTest() {
	jwtService := service.NewJWTService()
	s.jwtService = jwtService
}

func (s *DocumentEndpointsSuite) TestGetAllDocuments() {
	// Set up the environment
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	BASE_URL := os.Getenv("BASE_URL")
	if BASE_URL == "" {
		fmt.Println("BASE_URL environment variable is not set.")
		BASE_URL = "http://localhost:8080"
	}

	EMAIL := os.Getenv("TEST_USER_EMAIL")
	if EMAIL == "" {
		fmt.Println("TEST_USER_EMAIL environment variable is not set.")
		EMAIL = "khalid11abdu@gmail.com"
	}

	requestURL := BASE_URL + "/documents/getall"

	token := s.jwtService.GenerateToken(EMAIL, false)

	payload := map[string]string{"email": EMAIL}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatal(err)
	}

	// Create a new request to the server with the payload as the body
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.T().Fatal(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request to the server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	s.Assert().Equal(http.StatusOK, resp.StatusCode, "Unexpected status code")

	// Check the content type
	s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.T().Fatal(err)
	}

	// Optionally, parse the JSON response body for further assertions
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		s.T().Fatal(err)
	}
	documentsType := reflect.TypeOf(responseData["documents"])
	if documentsType != nil {
		s.Assert().Equal(documentsType.String(), "[]interface {}")
	}
}

func (s *DocumentEndpointsSuite) TestSearchDocuments() {
	// Set up the environment
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	BASE_URL := os.Getenv("BASE_URL")
	if BASE_URL == "" {
		fmt.Println("BASE_URL environment variable is not set.")
		BASE_URL = "http://localhost:8080"
	}

	EMAIL := os.Getenv("TEST_USER_EMAIL")
	if EMAIL == "" {
		fmt.Println("TEST_USER_EMAIL environment variable is not set.")
		EMAIL = "khalid11abdu@gmail.com"
	}

	SEARCH_QUERY := os.Getenv("SEARCH_QUERY")
	if SEARCH_QUERY == "" {
		fmt.Println("SEARCH_QUERY environment variable is not set.")
		SEARCH_QUERY = "Test"
	}

	requestURL := BASE_URL + "/documents/search"

	// Generate a valid JWT token for authentication
	token := s.jwtService.GenerateToken(EMAIL, false)

	// Prepare the payload for the request
	payload := map[string]string{
		"email":       EMAIL,
		"searchQuery": SEARCH_QUERY,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatal(err)
	}

	// Create a new request to the server with the payload as the body
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.T().Fatal(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request to the server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	s.Assert().Equal(http.StatusOK, resp.StatusCode, "Unexpected status code")

	// Check the content type
	s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.T().Fatal(err)
	}

	// Optionally, parse the JSON response body for further assertions
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		s.T().Fatal(err)
	}

	// check the type of the 'documents' field
	documentsType := reflect.TypeOf(responseData["documents"])
	if documentsType != nil {
		s.Assert().Equal(documentsType.String(), "[]interface {}")
	}
}

func (s *DocumentEndpointsSuite) TestCreateDocument() {
	// Set up the environment
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	BASE_URL := os.Getenv("BASE_URL")
	if BASE_URL == "" {
		fmt.Println("BASE_URL environment variable is not set.")
		BASE_URL = "http://localhost:8080"
	}

	EMAIL := os.Getenv("TEST_USER_EMAIL")
	if EMAIL == "" {
		fmt.Println("TEST_USER_EMAIL environment variable is not set.")
		EMAIL = "khalid11abdu@gmail.com"
	}

	TITLE := os.Getenv("TEST_DOCUMENT_TITLE")
	if TITLE == "" {
		fmt.Println("TEST_DOCUMENT_TITLE environment variable is not set.")
		TITLE = "Test Document"
	}

	BODY := os.Getenv("TEST_DOCUMENT_BODY")
	if BODY == "" {
		fmt.Println("TEST_DOCUMENT_BODY environment variable is not set.")
		BODY = "This is the body of the test document"
	}

	requestURL := BASE_URL + "/documents/createnew"

	// Generate a valid JWT token for authentication
	token := s.jwtService.GenerateToken(EMAIL, false)

	// Prepare the payload for the request
	payload := map[string]interface{}{
		"author":      EMAIL,
		"title":       TITLE,
		"body":        BODY,
		"readAccess":  []string{"readUser1@test.com", "readUser2@test.com"},
		"writeAccess": []string{"writeUser1@test.com", "writeUser2@test.com"},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatal(err)
	}

	// Create a new request to the server with the payload as the body
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.T().Fatal(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request to the server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	s.Assert().Equal(http.StatusCreated, resp.StatusCode, "Unexpected status code")

	// Check the content type
	s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.T().Fatal(err)
	}

	// Optionally, parse the JSON response body for further assertions
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		s.T().Fatal(err)
	}

	// check the type of the 'documentID' field
	documentIDType := reflect.TypeOf(responseData["document_id"])
	if documentIDType != nil {
		s.Assert().Equal(documentIDType.String(), "string")
	} else {
		s.T().Error("responseData['document_id'] is nil")
	}
}

func (s *DocumentEndpointsSuite) TestGetOneDocument() {
	// Set up the environment
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	BASE_URL := os.Getenv("BASE_URL")
	if BASE_URL == "" {
		fmt.Println("BASE_URL environment variable is not set.")
		BASE_URL = "http://localhost:8080"
	}

	EMAIL := os.Getenv("TEST_USER_EMAIL")
	if EMAIL == "" {
		fmt.Println("TEST_USER_EMAIL environment variable is not set.")
		EMAIL = "khalid11abdu@gmail.com"

		ID := os.Getenv("TEST_DOCUMENT_ID")
		if ID == "" {
			fmt.Println("TEST_DOCUMENT_ID environment variable is not set.")
			ID = "5f5e3e3e3e3e3e3e3e3e3e3e"
		}

		requestURL := BASE_URL + "/documents/getone/" + ID

		// Generate a valid JWT token for authentication
		token := s.jwtService.GenerateToken(EMAIL, true)

		// prepare a payload for the request
		payload := map[string]string{"id": ID, "email": EMAIL}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			s.T().Fatal(err)
		}

		// Create a new request to the server with the payload as the body
		req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			s.T().Fatal(err)
		}

		// Set the request headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		// Send the request to the server
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			s.T().Fatal(err)
		}

		// Check the HTTP response status code
		s.Assert().Equal(http.StatusOK, resp.StatusCode, "Unexpected status code")

		// Check the content type
		s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.T().Fatal(err)
		}

		// parse the JSON response body
		var responseData map[string]interface{}
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			s.T().Fatal(err)
		}

		// check the type of the 'document' field
		documentType := reflect.TypeOf(responseData["document"])
		if documentType != nil {
			s.Assert().Equal(documentType.String(), "map[string]interface {}")
		}
	}
}

func (s *DocumentEndpointsSuite) TestUpdateTitle() {
	// Set up the environment
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	BASE_URL := os.Getenv("BASE_URL")
	if BASE_URL == "" {
		fmt.Println("BASE_URL environment variable is not set.")
		BASE_URL = "http://localhost:8080"
	}

	EMAIL := os.Getenv("TEST_USER_EMAIL")
	if EMAIL == "" {
		fmt.Println("TEST_USER_EMAIL environment variable is not set.")
		EMAIL = "khalid11abdu@gmail.com"
	}

	ID := os.Getenv("TEST_DOCUMENT_ID")
	if ID == "" {
		fmt.Println("TEST_DOCUMENT_ID environment variable is not set.")
		ID = "5f5e3e3e3e3e3e3e3e3e3e3e"
	}

	UPDATED_TITLE := os.Getenv("TEST_DOCUMENT_UPDATED_TITLE")
	if UPDATED_TITLE == "" {
		fmt.Println("TEST_DOCUMENT_UPDATED_TITLE environment variable is not set.")
		UPDATED_TITLE = "Updated Test Document"
	}

	requestURL := BASE_URL + "/documents/updatetitle"

	// Generate a valid JWT token for authentication
	token := s.jwtService.GenerateToken(EMAIL, true)

	// Prepare the payload for the request
	payload := map[string]string{
		"id":           ID,
		"updatedTitle": UPDATED_TITLE,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatal(err)
	}

	// Create a new request to the server with the payload as the body
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.T().Fatal(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request to the server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

	// Check the HTTP response status code
	s.Assert().Equal(http.StatusOK, resp.StatusCode, "Unexpected status code")

	// Check the content type
	s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")
}

func (s *DocumentEndpointsSuite) TestUpdateCollaborators() {
	// Set up the environment
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	BASE_URL := os.Getenv("BASE_URL")
	if BASE_URL == "" {
		fmt.Println("BASE_URL environment variable is not set.")
		BASE_URL = "http://localhost:8080"
	}

	EMAIL := os.Getenv("TEST_USER_EMAIL")
	if EMAIL == "" {
		fmt.Println("TEST_USER_EMAIL environment variable is not set.")
		EMAIL = "khalid11abdu@gmail.com"
	}

	TITLE := os.Getenv("TEST_DOCUMENT_TITLE")
	if TITLE == "" {
		fmt.Println("TEST_DOCUMENT_TITLE environment variable is not set.")
		TITLE = "Test Document"
	}

	BODY := os.Getenv("TEST_DOCUMENT_BODY")
	if BODY == "" {
		fmt.Println("TEST_DOCUMENT_BODY environment variable is not set.")
		BODY = "This is the body of the test document"
	}

	requestURL := BASE_URL + "/documents/createnew"

	// Generate a valid JWT token for authentication
	token := s.jwtService.GenerateToken(EMAIL, false)

	// Prepare the payload for the request
	payload := map[string]interface{}{
		"author":      EMAIL,
		"title":       TITLE,
		"body":        BODY,
		"readAccess":  []string{"readUser1@test.com", "readUser2@test.com"},
		"writeAccess": []string{"writeUser1@test.com", "writeUser2@test.com"},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatal(err)
	}

	// Create a new request to the server with the payload as the body
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.T().Fatal(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request to the server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	s.Assert().Equal(http.StatusCreated, resp.StatusCode, "Unexpected status code")

	// Check the content type
	s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.T().Fatal(err)
	}

	// Optionally, parse the JSON response body for further assertions
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		s.T().Fatal(err)
	}

	// check the type of the 'documentID' field
	documentIDType := reflect.TypeOf(responseData["document_id"])
	if documentIDType != nil {
		s.Assert().Equal(documentIDType.String(), "string")
	} else {
		s.T().Error("responseData['document_id'] is nil")
	}

	ID := fmt.Sprintf("%v", responseData["document_id"])

	requestURL = BASE_URL + "/documents/updatecollaborators"

	// Prepare the payload for the request
	payload = map[string]interface{}{
		"document_id": ID,
		"readAccess":  []string{"test@example.com"},
		"writeAccess": []string{"test2@example.com"},
	}

	payloadBytes, err = json.Marshal(payload)
	if err != nil {
		s.T().Fatal(err)
	}

	// Create a new request to the server with the payload as the body
	req, err = http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.T().Fatal(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request to the server
	resp, err = client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

	// Check the HTTP response status code
	s.Assert().Equal(http.StatusOK, resp.StatusCode, "Unexpected status code")

	// Check the content type
	s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")
}

func (s *DocumentEndpointsSuite) TestDeleteDocument() {
	// Set up the environment
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	BASE_URL := os.Getenv("BASE_URL")
	if BASE_URL == "" {
		fmt.Println("BASE_URL environment variable is not set.")
		BASE_URL = "http://localhost:8080"
	}

	EMAIL := os.Getenv("TEST_USER_EMAIL")
	if EMAIL == "" {
		fmt.Println("TEST_USER_EMAIL environment variable is not set.")
		EMAIL = "khalid11abdu@gmail.com"
	}

	TITLE := os.Getenv("TEST_DOCUMENT_TITLE")
	if TITLE == "" {
		fmt.Println("TEST_DOCUMENT_TITLE environment variable is not set.")
		TITLE = "Test Document"
	}

	BODY := os.Getenv("TEST_DOCUMENT_BODY")
	if BODY == "" {
		fmt.Println("TEST_DOCUMENT_BODY environment variable is not set.")
		BODY = "This is the body of the test document"
	}

	requestURL := BASE_URL + "/documents/createnew"

	// Generate a valid JWT token for authentication
	token := s.jwtService.GenerateToken(EMAIL, false)

	// Prepare the payload for the request
	payload := map[string]interface{}{
		"author":      EMAIL,
		"title":       TITLE,
		"body":        BODY,
		"readAccess":  []string{"readUser1@test.com", "readUser2@test.com"},
		"writeAccess": []string{"writeUser1@test.com", "writeUser2@test.com"},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatal(err)
	}

	// Create a new request to the server with the payload as the body
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.T().Fatal(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request to the server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

	// Check the HTTP response status code
	s.Assert().Equal(http.StatusCreated, resp.StatusCode, "Unexpected status code")

	// Check the content type
	s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.T().Fatal(err)
	}

	// parse the JSON response body for further assertions
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		s.T().Fatal(err)
	}

	// check the type of the 'documentID' field
	documentIDType := reflect.TypeOf(responseData["document_id"])
	if documentIDType != nil {
		s.Assert().Equal(documentIDType.String(), "string")
	} else {
		s.T().Error("responseData['document_id'] is nil")
	}

	requestURL = BASE_URL + "/documents/delete/" + fmt.Sprintf("%v", responseData["document_id"])

	// Create a new request to the server
	req, err = http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		s.T().Fatal(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request to the server
	resp, err = client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

	// Check the HTTP response status code
	s.Assert().Equal(http.StatusOK, resp.StatusCode, "Unexpected status code")

	// Check the content type
	s.Assert().Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected content type")
}
