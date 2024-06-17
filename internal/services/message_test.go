package services

import (
	"chat-system/internal/models"
	"chat-system/mocks"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Unit tests

type MessageUnitTestSuite struct {
	suite.Suite
	msgService *mocks.MessageService
}

func TestMessageUnitTestSuite(t *testing.T) {
	suite.Run(t, new(MessageUnitTestSuite))
}

func (uts *MessageUnitTestSuite) SetupTest() {
	uts.msgService = new(mocks.MessageService)
}

func (uts *MessageUnitTestSuite) TestCreateMessage() {
	uts.msgService.On("CreateMessage", mock.Anything).Return(nil)

	msg := &models.Message{
		Sender:    "a",
		Recipient: "b",
		Content:   "This is a content",
	}

	err := uts.msgService.CreateMessage(msg)

	uts.Nil(err)
}

func (uts *MessageUnitTestSuite) TestGetMessages() {
	testUsername := "test"
	testMsg := models.Message{
		ID:        gocql.TimeUUID(),
		Sender:    "a",
		Recipient: "b",
		Timestamp: time.Now(),
		Content:   "This is a test content.",
	}
	expectedMsgs := []models.Message{}
	expectedMsgs = append(expectedMsgs, testMsg)

	uts.msgService.On("GetMessages", mock.Anything).Return(expectedMsgs, nil)

	actual, err := uts.msgService.GetMessages(testUsername)

	uts.Nil(err)
	uts.Equal(expectedMsgs, actual)
	uts.Contains(expectedMsgs, actual[0])
}
