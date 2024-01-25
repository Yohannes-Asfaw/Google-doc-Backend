package unit_tests

import (
	"testing"

	"github.com/khallihub/godoc/service"
	"github.com/stretchr/testify/suite"
)

type JWTServiceTestSuite struct {
	suite.Suite
	jwtService     service.JWTService
}

func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, &JWTServiceTestSuite{})
}

func (uts *JWTServiceTestSuite) SetupTest() {
	jwtService := service.NewJWTService()
	uts.jwtService = jwtService
}

func (uts *JWTServiceTestSuite) TestGenerateToken() {
	// Mock data
	name := "author@test.com"
	admin := true

	// Call the GenerateToken method with the mocked data
	token := uts.jwtService.GenerateToken(name, admin)

	// Assert that token was generated
	uts.NotEmpty(token)
}

func (uts *JWTServiceTestSuite) TestValidateToken() {
	// Mock data
	name := "test"
	admin := true

	// Call the GenerateToken method with the mocked data
	token := uts.jwtService.GenerateToken(name, admin)

	// Call the ValidateToken method with the generated token
	_, err := uts.jwtService.ValidateToken(token)

	// Assert that token was validated
	uts.Nil(err)
}