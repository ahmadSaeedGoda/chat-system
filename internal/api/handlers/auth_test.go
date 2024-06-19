package handlers

import (
	"bytes"
	common "chat-system/internal/api/common/constants"
	"chat-system/internal/api/middlewares"
	"chat-system/internal/api/transformers"
	"chat-system/internal/models"
	"chat-system/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	handler *userHandler
	service *mocks.UserService
	server  *httptest.Server
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (ats *AuthTestSuite) SetupTest() {
	r := mux.NewRouter()
	r.Use(middlewares.HandleErrors)

	ats.service = &mocks.UserService{}
	ats.handler = NewUserHandler(ats.service)

	r.HandleFunc("/register", ats.handler.Register).Methods("POST")
	r.HandleFunc("/login", ats.handler.Login).Methods("POST")

	ats.server = httptest.NewServer(r)
}

func (ats *AuthTestSuite) TearDownTest() {
	ats.server.Close()
}

func (ats *AuthTestSuite) TestRegister_Invalid_Input() {
	// Test invalid req body no username
	testPassword := "123456"
	invalidBody := &models.RegisterInput{Password: testPassword}
	body, err := json.Marshal(invalidBody)
	ats.NoError(err, "Failed to marshal registerInput")
	expectedUser := &models.User{Username: "user1", Password: "123456"}
	ats.service.On("UserExists", mock.Anything).Return(false, nil).Once()
	ats.service.On("CreateUser", mock.Anything).Return(expectedUser, nil).Once()
	resp, err := http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")

	ats.Equal(http.StatusBadRequest, resp.StatusCode)

	// Test invalid req body no password
	testUserName := "user1"
	invalidBody = &models.RegisterInput{Username: testUserName}
	body, err = json.Marshal(invalidBody)
	ats.NoError(err, "Failed to marshal registerInput")
	expectedUser = &models.User{Username: "user1", Password: "123456"}
	ats.service.On("UserExists", mock.Anything).Return(false, nil).Once()
	ats.service.On("CreateUser", mock.Anything).Return(expectedUser, nil).Once()
	resp, err = http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")

	ats.Equal(http.StatusBadRequest, resp.StatusCode)

	var response = struct {
		Error string `json:"error"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	ats.NoError(err, "Failed to decode response body")

	ats.Equal(common.BAD_REQUEST, response.Error)

	// Test invalid username min
	testPassword = "123456"
	registerInput := &models.RegisterInput{Username: "", Password: testPassword}
	body, err = json.Marshal(registerInput)
	ats.NoError(err, "Failed to marshal registerInput")

	resp, err = http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")

	ats.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	ats.NoError(err, "Failed to decode response body")

	ats.Equal(common.BAD_REQUEST, response.Error)

	// test invalid username max
	testUserName = "ThisIsMoreThanSixteenCharLength"
	registerInput = &models.RegisterInput{Username: testUserName, Password: testPassword}
	body, err = json.Marshal(registerInput)
	ats.NoError(err, "Failed to marshal registerInput")

	resp, err = http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")
	defer resp.Body.Close()

	ats.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	ats.NoError(err, "Failed to decode response body")

	ats.Equal(common.BAD_REQUEST, response.Error)

	// test invalid password min
	testUserName = "user1"
	testPassword = "12345"
	registerInput = &models.RegisterInput{Username: testUserName, Password: testPassword}
	body, err = json.Marshal(registerInput)
	ats.NoError(err, "Failed to marshal registerInput")

	resp, err = http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")
	defer resp.Body.Close()

	ats.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	ats.NoError(err, "Failed to decode response body")

	ats.Equal(common.BAD_REQUEST, response.Error)

	// test invalid password max
	testPassword = "ThisIsMoreThanThirteenCharLength"
	registerInput = &models.RegisterInput{Username: testUserName, Password: testPassword}
	body, err = json.Marshal(registerInput)
	ats.NoError(err, "Failed to marshal registerInput")

	resp, err = http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")
	defer resp.Body.Close()

	ats.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	ats.NoError(err, "Failed to decode response body")

	ats.Equal(common.BAD_REQUEST, response.Error)

}

func (ats *AuthTestSuite) TestRegister_Username_Taken() {
	ats.service.On("UserExists", mock.Anything).Return(true, nil).Once()

	testUserName := "user1"
	registerInput := &models.RegisterInput{Username: testUserName, Password: "123456"}
	body, err := json.Marshal(registerInput)
	ats.NoError(err, "Failed to marshal registerInput")

	resp, err := http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")
	defer resp.Body.Close()

	ats.Equal(http.StatusConflict, resp.StatusCode)

	var response = struct {
		Error string `json:"error"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	ats.NoError(err, "Failed to decode response body")

	ats.Equal(common.REGISTER_USER_EXISTS, response.Error, "Expected error message for existing user")
}

func (ats *AuthTestSuite) TestRegister_Unexpected_Err() {
	expectedErr := errors.New("Internal Server Error")
	ats.service.On("UserExists", mock.Anything).Return(false, expectedErr).Once()

	testUserName := "user1"
	testPassword := "123456"
	registerInput := &models.RegisterInput{Username: testUserName, Password: testPassword}
	body, err := json.Marshal(registerInput)
	ats.NoError(err, "Failed to marshal registerInput")

	resp, err := http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")
	defer resp.Body.Close()

	ats.Equal(http.StatusInternalServerError, resp.StatusCode)

	var response = struct {
		Error string `json:"error"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	ats.NoError(err, "Failed to decode response body")

	ats.Equal(expectedErr.Error(), response.Error)
}

func (ats *AuthTestSuite) TestRegister_Success() {
	ats.service.On("UserExists", mock.Anything).Return(false, nil).Once()

	testUsername, testPassword := "user1", "123456"
	expectedUser := &models.User{Username: testUsername, Password: testPassword}
	ats.service.On("CreateUser", mock.Anything).Return(expectedUser, nil).Once()

	registerInput := models.RegisterInput{Username: testUsername, Password: testPassword}
	body, err := json.Marshal(registerInput)
	ats.NoError(err, "Failed to marshal registerInput")

	resp, err := http.Post(ats.server.URL+"/register", "application/json", bytes.NewBuffer(body))
	ats.NoError(err, "Failed to make POST request")
	defer resp.Body.Close()

	ats.Equal(http.StatusCreated, resp.StatusCode)

	var ur transformers.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&ur)
	ats.NoError(err, "Failed to decode response body")

	ats.Equal(testUsername, ur.Username)
}
