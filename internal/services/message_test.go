package services

import (
	"chat-system/internal/models"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/suite"
)

type MessagesTestSuite struct {
	suite.Suite
	dbSession *gocql.Session
	tableName string
	service   *messageService
}

func (mts *MessagesTestSuite) DBSession() *gocql.Session {
	return mts.dbSession
}

func (mts *MessagesTestSuite) SetDBSession(session *gocql.Session) {
	mts.dbSession = session
}

func (mts *MessagesTestSuite) DBTable() string {
	return mts.tableName
}

func (mts *MessagesTestSuite) SetDBTable(tableName string) {
	mts.tableName = tableName
}

func (mts *MessagesTestSuite) Service() *messageService {
	return mts.service
}

func (mts *MessagesTestSuite) SetService(service *messageService) {
	mts.service = service
}

func TestMessagesTestSuite(t *testing.T) {
	suite.Run(t, new(MessagesTestSuite))
}

func (mts *MessagesTestSuite) SetupSuite() {
	mts.SetDBTable(MSGS_TEST_TABLE_NAME)
	mts.SetDBSession(setupDatabase(mts))

	mts.SetService(NewMessageService(mts.DBSession(), KEYSPACE_TEST, mts.DBTable()))
}

func (mts *MessagesTestSuite) TearDownSuite() {
	tearDownDatabase(mts)
}

func (mts *MessagesTestSuite) TestCreateMessage_Success() {
	msg := &models.Message{
		ID:        gocql.TimeUUID(),
		Sender:    "a",
		Recipient: "b",
		Timestamp: time.Now(),
		Content:   "test",
	}

	err := mts.Service().CreateMessage(msg)

	mts.Nil(err)

	cleanTable(mts)
}

func (mts *MessagesTestSuite) TestCreateMessage_ErrDB() {
	errMsg := "Key may not be empty"

	err := mts.Service().CreateMessage(&models.Message{})

	mts.EqualError(err, errMsg)
}

func (mts *MessagesTestSuite) TestGetMessages_Success() {
	cleanTable(mts)

	seedTable(mts)

	testUsername := "user1"
	actual, err := mts.Service().GetMessages(testUsername)

	mts.Nil(err)
	mts.NotEmpty(actual)
	mts.Len(actual, 1)
	mts.Equal(actual[0].Content, "test content")

	cleanTable(mts)
}

func (mts *MessagesTestSuite) TestGetMessages_Empty() {
	testUsername := "test-non-exist"

	actual, err := mts.Service().GetMessages(testUsername)

	mts.Nil(err)
	mts.Empty(actual)
}

func (mts *MessagesTestSuite) TestGetMessages_ErrDB() {
	errMsg := "Key may not be empty"

	actual, err := mts.Service().GetMessages("")

	mts.EqualError(err, errMsg)
	mts.Nil(actual)
}
