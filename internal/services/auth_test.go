package services

import (
	"chat-system/internal/models"
	"chat-system/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Unit tests

type AuthUnitTestSuite struct {
	suite.Suite
	userService *mocks.UserService
}

func TestAuthUnitTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUnitTestSuite))
}

func (uts *AuthUnitTestSuite) SetupTest() {
	uts.userService = new(mocks.UserService)
}

func (uts *AuthUnitTestSuite) TestUserExists() {
	uts.userService.On("UserExists", mock.Anything).Return(true, nil)

	exists, err := uts.userService.UserExists("user1")

	uts.Nil(err)
	uts.True(exists)
}

func (uts *AuthUnitTestSuite) TestCreateUser() {
	userInput := &models.RegisterInput{Username: "user1", Password: "password123"}
	expectedUser := &models.User{Username: "user1"}

	uts.userService.On("CreateUser", mock.Anything).Return(expectedUser, nil)

	actual, err := uts.userService.CreateUser(userInput)

	uts.Nil(err)
	uts.Equal(expectedUser, actual)
}
