package services

import (
	"chat-system/internal/models"
	"errors"
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	dbSession *gocql.Session
	tableName string
	service   *userService
}

func (ats *AuthTestSuite) DBSession() *gocql.Session {
	return ats.dbSession
}

func (ats *AuthTestSuite) SetDBSession(session *gocql.Session) {
	ats.dbSession = session
}

func (ats *AuthTestSuite) DBTable() string {
	return ats.tableName
}

func (ats *AuthTestSuite) SetDBTable(tableName string) {
	ats.tableName = tableName
}

func (ats *AuthTestSuite) Service() *userService {
	return ats.service
}

func (ats *AuthTestSuite) SetService(service *userService) {
	ats.service = service
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (ats *AuthTestSuite) SetupSuite() {
	ats.SetDBTable(USERS_TEST_TABLE_NAME)
	ats.SetDBSession(setupDatabase(ats))

	ats.SetService(NewUserService(ats.DBSession(), KEYSPACE_TEST, ats.DBTable()))
}

func (ats *AuthTestSuite) TearDownSuite() {
	tearDownDatabase(ats)
}

func (ats *AuthTestSuite) TestUserExists_True() {
	seedTable(ats)

	testUsername := "user1"
	actual, err := ats.Service().UserExists(testUsername)

	ats.Nil(err)
	ats.True(actual)
}

func (ats *AuthTestSuite) TestUserExists_False() {
	testUsername := "non-exist-for-sure"
	actual, err := ats.Service().UserExists(testUsername)

	ats.Nil(err)
	ats.False(actual)
}

func (ats *AuthTestSuite) TestCreateUser_Success() {
	cleanTable(ats)

	userInput := &models.RegisterInput{Username: "user1", Password: "password123"}
	expectedUser := &models.User{Username: "user1"}

	actual, err := ats.Service().CreateUser(userInput)

	ats.Nil(err)
	ats.Equal(expectedUser.Username, actual.Username)
}

func (ats *AuthTestSuite) TestCreateUser_ErrorHash() {
	errMsg := "bcrypt: password length exceeds 72 bytes"
	expectedErr := errors.New(errMsg)

	veryLongPass := "This is a test string that has exactly seventy-three characters in length"

	userInput := &models.RegisterInput{Username: "user1", Password: veryLongPass}

	actual, err := ats.Service().CreateUser(userInput)

	ats.Nil(actual)
	ats.Equal(expectedErr, err)
	ats.EqualError(expectedErr, errMsg)
}
