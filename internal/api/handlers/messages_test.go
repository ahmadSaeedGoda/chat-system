package handlers

import (
	"bytes"
	common "chat-system/internal/api/common/constants"
	"chat-system/internal/api/middlewares"
	"chat-system/internal/models"
	"chat-system/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type ErrRes struct {
	Error string `json:"error"`
}

type MessagesTestSuite struct {
	suite.Suite
	handler                MsgHandler
	msgService             *mocks.MessageService
	userService            *mocks.UserService
	server                 *httptest.Server
	reqContentType         string
	sendEndpointReqPath    string
	getMsgsEndpointReqPath string
	errResponse            *ErrRes
}

func TestMessagesTestSuite(t *testing.T) {
	suite.Run(t, new(MessagesTestSuite))
}

func (mts *MessagesTestSuite) SetupSuite() {
	mts.reqContentType = "application/json"
	mts.sendEndpointReqPath = "/send"
	mts.getMsgsEndpointReqPath = "/messages"
	mts.errResponse = &ErrRes{}
}

func (mts *MessagesTestSuite) SetupTest() {
	r := mux.NewRouter()
	r.Use(middlewares.HandleErrors)

	mts.msgService = &mocks.MessageService{}
	mts.userService = &mocks.UserService{}
	mts.handler = NewMsgHandler(mts.msgService, mts.userService)

	r.HandleFunc("/send", mts.handler.SendMessage).Methods("POST")
	r.HandleFunc("/messages", mts.handler.GetMessages).Methods("GET")

	mts.server = httptest.NewServer(r)
}

func (mts *MessagesTestSuite) TearDownTest() {
	mts.server.Close()
}

func (mts *MessagesTestSuite) Test_Send_Invalid_Input_NoRecipient() {
	content := "This is a test content"
	invalidBody := &models.SendMessageInput{Content: content}
	body, err := json.Marshal(invalidBody)
	mts.NoError(err, "Failed to marshal registerInput")
	expectedErr := errors.New(common.BAD_REQUEST)
	resp, err := http.Post(mts.server.URL+mts.sendEndpointReqPath, mts.reqContentType, bytes.NewBuffer(body))
	mts.NoError(err, "Failed to make POST request")

	mts.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&mts.errResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(common.BAD_REQUEST, mts.errResponse.Error)
	mts.Equal(expectedErr.Error(), mts.errResponse.Error)
}

func (mts *MessagesTestSuite) Test_Send_Invalid_Input_NoContent() {
	recipient := "x"
	invalidBody := &models.SendMessageInput{Recipient: recipient}
	body, err := json.Marshal(invalidBody)
	mts.NoError(err, "Failed to marshal registerInput")

	expectedErr := errors.New(common.BAD_REQUEST)
	resp, err := http.Post(mts.server.URL+mts.sendEndpointReqPath, mts.reqContentType, bytes.NewBuffer(body))
	mts.NoError(err, "Failed to make POST request")

	mts.Equal(http.StatusBadRequest, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&mts.errResponse)
	mts.NoError(err, "Failed to decode response body")

	mts.Equal(common.BAD_REQUEST, mts.errResponse.Error)
	mts.Equal(expectedErr.Error(), mts.errResponse.Error)
}
