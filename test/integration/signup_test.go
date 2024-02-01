package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type SignupEndpointsSuite struct {
	suite.Suite
}

func TestSignupEndpointsSuite(t *testing.T) {
	suite.Run(t, new(SignupEndpointsSuite))
}

func (s *SignupEndpointsSuite) TestSignup() {
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

	EMAIL := os.Getenv("TEST_NEW_USER_EMAIL")
	if EMAIL == "" {
		fmt.Println("TEST_USER_EMAIL environment variable is not set.")
		EMAIL = "khalid11abdu@gmail.com"
	}

	PASSWORD := os.Getenv("TEST_NEW_USER_PASSWORD")
	if PASSWORD == "" {
		fmt.Println("TEST_USER_PASSWORD environment variable is not set.")
		PASSWORD = "Khalid1!"
	}

	requestURL := BASE_URL + "/auth/signup"

	// Prepare the payload for the request
	payload := map[string]string{
		"email":    EMAIL,
		"password": PASSWORD,
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

	// parse the JSON response body for further assertions
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		s.T().Fatal(err)
	}

	// Check the presence of the token in the response
	message := responseData["message"].(string)
	s.Assert().Equal("User created successfully", message)
}	
