package handlers

import (
	"bytes"
	"chat-system/internal/api/auth"
	common "chat-system/internal/api/common/constants"
	"chat-system/internal/api/common/responses"
	"chat-system/internal/api/middlewares"
	"chat-system/internal/models"
	"chat-system/mocks"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MessagesTestSuite struct {
	suite.Suite
	msgService         *mocks.MessageService
	userService        *mocks.UserService
	sendEndpointUrl    string
	getMsgsEndpointUrl string
	authHeader         string
	handler            *msgHandler
	errResponse        *responses.ErrResponse
	middleware         func(handlerMethod func(w http.ResponseWriter, r *http.Request)) http.Handler
}

func TestMessagesTestSuite(t *testing.T) {
	suite.Run(t, new(MessagesTestSuite))
}

func (mts *MessagesTestSuite) SetupSuite() {
	// Set up environment variable
	os.Setenv("AUTH_HEADER_PREFIX", "Bearer")
	os.Setenv("JWT_SECRET_KEY", "secret")

	mts.msgService = &mocks.MessageService{}
	mts.userService = &mocks.UserService{}

	reqSenderUsername := "User1"
	token, err := auth.GenerateToken(reqSenderUsername)
	mts.NoError(err, "Failed to create token")

	mts.authHeader = fmt.Sprintf("Bearer %s", token)

	mts.handler = NewMsgHandler(mts.msgService, mts.userService)

	mts.middleware = func(handlerMethod func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return middlewares.IsAuth(
			http.HandlerFunc(
				middlewares.HandleErrors(
					http.HandlerFunc(handlerMethod),
				).ServeHTTP,
			),
		)
	}

	mts.sendEndpointUrl = "localhost/send"
	mts.getMsgsEndpointUrl = "localhost/api/v1/messages"

	mts.errResponse = &responses.ErrResponse{}
}

func (mts *MessagesTestSuite) Test_Send_Invalid_Input_NoRecipient() {
	content := "This is a test content"
	invalidBody := &models.SendMessageInput{Content: content}
	body, err := json.Marshal(invalidBody)
	mts.NoError(err, "Failed to marshal registerInput")

	req, err := http.NewRequest("POST", mts.sendEndpointUrl, bytes.NewBuffer(body))
	mts.NoError(err, "Failed to make POST request")

	req.Header.Set("Authorization", mts.authHeader)

	rr := httptest.NewRecorder()

	mts.middleware(mts.handler.SendMessage).ServeHTTP(rr, req)

	resp := rr.Result()

	mts.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&mts.errResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(common.BAD_REQUEST, mts.errResponse.Error)
}

func (mts *MessagesTestSuite) Test_Send_Invalid_Input_NoContent() {
	recipient := "x"
	invalidBody := &models.SendMessageInput{Recipient: recipient}
	body, err := json.Marshal(invalidBody)
	mts.NoError(err, "Failed to marshal registerInput")

	req, err := http.NewRequest("POST", mts.sendEndpointUrl, bytes.NewBuffer(body))
	mts.NoError(err, "Failed to make POST request")

	req.Header.Set("Authorization", mts.authHeader)

	rr := httptest.NewRecorder()

	mts.middleware(mts.handler.SendMessage).ServeHTTP(rr, req)

	resp := rr.Result()

	mts.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&mts.errResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(common.BAD_REQUEST, mts.errResponse.Error)
}

func (mts *MessagesTestSuite) Test_Send_Err_Checking_Recipient() {
	expectedErr := errors.New(common.INTERNAL_SERVER_ERROR)
	mts.userService.On("UserExists", mock.Anything).Return(false, expectedErr).Once()

	recipient := "x"
	content := "Test"
	reqBody := &models.SendMessageInput{Recipient: recipient, Content: content}
	body, err := json.Marshal(reqBody)
	mts.NoError(err, "Failed to marshal registerInput")

	req, err := http.NewRequest("POST", mts.sendEndpointUrl, bytes.NewBuffer(body))
	mts.NoError(err, "Failed to make POST request")

	req.Header.Set("Authorization", mts.authHeader)

	rr := httptest.NewRecorder()

	mts.middleware(mts.handler.SendMessage).ServeHTTP(rr, req)

	resp := rr.Result()

	mts.Equal(http.StatusInternalServerError, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&mts.errResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(common.INTERNAL_SERVER_ERROR, mts.errResponse.Error)
}

func (mts *MessagesTestSuite) Test_Send_Recipient_Not_Found() {
	mts.userService.On("UserExists", mock.Anything).Return(false, nil).Once()

	recipient := "x"
	content := "Test"
	reqBody := &models.SendMessageInput{Recipient: recipient, Content: content}
	body, err := json.Marshal(reqBody)
	mts.NoError(err, "Failed to marshal registerInput")

	req, err := http.NewRequest("POST", mts.sendEndpointUrl, bytes.NewBuffer(body))
	mts.NoError(err, "Failed to make POST request")

	req.Header.Set("Authorization", mts.authHeader)

	rr := httptest.NewRecorder()

	mts.middleware(mts.handler.SendMessage).ServeHTTP(rr, req)

	resp := rr.Result()

	mts.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&mts.errResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(common.SEND_MESSAGE_NO_RECIPIENT, mts.errResponse.Error)
}

func (mts *MessagesTestSuite) Test_Send_DB_Err() {
	mts.userService.On("UserExists", mock.Anything).Return(true, nil).Once()

	expectedErr := errors.New(common.INTERNAL_SERVER_ERROR)
	mts.msgService.On("CreateMessage", mock.Anything).Return(expectedErr).Once()

	recipient := "x"
	content := "Test"
	reqBody := &models.SendMessageInput{Recipient: recipient, Content: content}
	body, err := json.Marshal(reqBody)
	mts.NoError(err, "Failed to marshal registerInput")

	req, err := http.NewRequest("POST", mts.sendEndpointUrl, bytes.NewBuffer(body))
	mts.NoError(err, "Failed to make POST request")

	req.Header.Set("Authorization", mts.authHeader)

	rr := httptest.NewRecorder()

	mts.middleware(mts.handler.SendMessage).ServeHTTP(rr, req)

	resp := rr.Result()

	mts.Equal(http.StatusInternalServerError, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&mts.errResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(common.INTERNAL_SERVER_ERROR, mts.errResponse.Error)
	mts.Equal(expectedErr.Error(), mts.errResponse.Error)
}

func (mts *MessagesTestSuite) Test_Send_Success() {
	recipient := "User2"
	content := "Test Content"
	reqBody := &models.SendMessageInput{Recipient: recipient, Content: content}
	body, err := json.Marshal(reqBody)
	mts.NoError(err, "Failed to marshal registerInput")

	req, err := http.NewRequest("POST", mts.sendEndpointUrl, bytes.NewBuffer(body))
	mts.NoError(err, "Failed to make POST request")

	req.Header.Set("Authorization", mts.authHeader)

	mts.userService.On("UserExists", mock.Anything).Return(true, nil).Once()
	mts.msgService.On("CreateMessage", mock.Anything).Return(nil).Once()
	mts.msgService.On("UpdateCachedMsgsForUser", mock.Anything, mock.Anything).Return(nil).Twice()

	rr := httptest.NewRecorder()

	mts.middleware(mts.handler.SendMessage).ServeHTTP(rr, req)

	resp := rr.Result()

	mts.Equal(http.StatusCreated, resp.StatusCode)

	msgResponse := &responses.MessageResponse{}
	err = json.NewDecoder(resp.Body).Decode(&msgResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(msgResponse.Recipient, recipient)
	mts.Equal(msgResponse.Content, content)
	mts.Equal(msgResponse.Sender, "User1")
}

func (mts *MessagesTestSuite) Test_GetMessages_Success() {
	req, err := http.NewRequest("GET", mts.getMsgsEndpointUrl, nil)
	mts.NoError(err, "Failed to make request")

	req.Header.Set("Authorization", mts.authHeader)

	mts.msgService.On("GetFromCache", mock.Anything).Return(nil, nil).Once()

	expectedMsg := models.Message{Sender: "Mickey", Recipient: "Minnie", Content: "Hi"}
	msgsArr := make([]models.Message, 0)
	msgsArr = append(msgsArr, expectedMsg)
	mts.msgService.On("GetMessages", mock.Anything).Return(msgsArr, nil).Once()

	mts.msgService.On("SetMessagesToCache", mock.Anything, mock.Anything).Return(nil).Once()

	rr := httptest.NewRecorder()

	mts.middleware(mts.handler.GetMessages).ServeHTTP(rr, req)

	resp := rr.Result()

	mts.Equal(http.StatusOK, resp.StatusCode)

	var msgsResponse responses.MessagesResponse
	err = json.NewDecoder(resp.Body).Decode(&msgsResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Len(msgsResponse.Messages, 1)
	mts.Equal(msgsResponse.Messages[0].Sender, expectedMsg.Sender)
	mts.Equal(msgsResponse.Messages[0].Recipient, expectedMsg.Recipient)
	mts.Equal(msgsResponse.Messages[0].Content, expectedMsg.Content)
}

func (mts *MessagesTestSuite) Test_GetMessages_Error() {
	req, err := http.NewRequest("GET", mts.getMsgsEndpointUrl, nil)
	mts.NoError(err, "Failed to make request")

	req.Header.Set("Authorization", mts.authHeader)

	mts.msgService.On("GetFromCache", mock.Anything).Return(nil, nil).Once()

	expectedErr := errors.New("DB is Down :(")
	mts.msgService.On("GetMessages", mock.Anything).Return(nil, expectedErr).Once()

	mts.msgService.On("SetMessagesToCache", mock.Anything, mock.Anything).Return(nil).Once()

	rr := httptest.NewRecorder()

	mts.middleware(mts.handler.GetMessages).ServeHTTP(rr, req)

	resp := rr.Result()

	mts.Equal(http.StatusInternalServerError, resp.StatusCode)

	var errResponse responses.ErrResponse
	err = json.NewDecoder(resp.Body).Decode(&errResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(errResponse.Error, common.INTERNAL_SERVER_ERROR)
}
