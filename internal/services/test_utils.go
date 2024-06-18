// Helper functions & consts for running tests

package services

import (
	"chat-system/internal/cassandra"
	"fmt"
	"testing"
	"time"

	"github.com/gocql/gocql"
)

const (
	KEYSPACE_TEST         = "chat_test"
	REPLICA_COUNT         = 2
	CLUSTER_ADDRS         = "127.0.0.1"
	MSGS_TEST_TABLE_NAME  = "messages"
	USERS_TEST_TABLE_NAME = "users"
)

type TestSuite interface {
	T() *testing.T
	FailNowf(failureMessage string, msg string, args ...interface{}) bool
	// Get Current DB Session Obj
	DBSession() *gocql.Session
	// Set Current DB Session Obj
	SetDBSession(session *gocql.Session)
	// Get DB Table Name
	DBTable() string
	// Set DB Table Name
	SetDBTable(tableName string)
}

func setupDatabase(ts TestSuite) *gocql.Session {
	// Establish Conn
	ts.T().Log("setting up test database")

	dbSession, err := cassandra.Init(KEYSPACE_TEST, REPLICA_COUNT, CLUSTER_ADDRS)
	if err != nil {
		ts.FailNowf("unable to connect to test database", err.Error())
	}
	ts.T().Log("connected to test database")

	// Configure test suite
	ts.SetDBSession(dbSession)

	// Create test table
	err = ts.DBSession().Query(getCreateTestTableStmt(ts)).Exec()
	if err != nil {
		ts.FailNowf("unable to create table", err.Error())
	}
	ts.T().Logf("created '%s' table in test database", ts.DBTable())
	return dbSession
}

func getCreateTestTableStmt(ts TestSuite) string {
	var createTableStmt string

	switch ts.DBTable() {
	case MSGS_TEST_TABLE_NAME:
		createTableStmt = fmt.Sprintf(
			`
				CREATE TABLE IF NOT EXISTS %s.%s
				(
					user TEXT,
					timestamp TIMESTAMP,
					id UUID,
					sender TEXT,
					recipient TEXT,
					content TEXT,
					PRIMARY KEY (user, timestamp, id)
				)
				WITH CLUSTERING ORDER BY (timestamp DESC)
			`,
			KEYSPACE_TEST,
			ts.DBTable(),
		)
	case USERS_TEST_TABLE_NAME:
		createTableStmt = fmt.Sprintf(
			`
				CREATE TABLE IF NOT EXISTS %s.%s (
					id UUID,
					username TEXT PRIMARY KEY,
					password TEXT,
				)

			`,
			KEYSPACE_TEST,
			ts.DBTable(),
		)
	}

	return createTableStmt
}

func seedTable(ts TestSuite) {
	ts.T().Logf("seeding test table for test '%s'", ts.T().Name())

	userPrefix := "user"
	var errTable, errRecord error
	for i := 1; i <= 2; i++ {
		userSuffix := fmt.Sprint(i)
		switch ts.DBTable() {
		case MSGS_TEST_TABLE_NAME:
			errRecord = ts.DBSession().Query(
				getSeedStmt(ts),
				userPrefix+userSuffix,
				time.Now(),
				gocql.TimeUUID(),
				userPrefix+userSuffix,
				userPrefix+fmt.Sprint(i+1),
				"test content",
			).Exec()

		case USERS_TEST_TABLE_NAME:
			errRecord = ts.DBSession().Query(
				getSeedStmt(ts),
				userPrefix+userSuffix,
				"",
			).Exec()
		}
		if errRecord != nil {
			ts.FailNowf(
				"unable to seed one record for test table",
				errTable.Error(),
			)
		}
	}

	if errTable != nil {
		ts.FailNowf("unable to seed test table", errTable.Error())
	}

	ts.T().Logf("db seeded successfully for test '%s'", ts.T().Name())
}

func getSeedStmt(ts TestSuite) string {
	switch ts.DBTable() {
	case MSGS_TEST_TABLE_NAME:
		return fmt.Sprintf(
			`INSERT INTO %s.%s
			(user, timestamp, id, sender, recipient, content)
			VALUES (?, ?, ?, ?, ?, ?)`,
			KEYSPACE_TEST,
			ts.DBTable(),
		)
	case USERS_TEST_TABLE_NAME:
		return fmt.Sprintf(
			`INSERT INTO %s.%s
			(username, password)
			VALUES (?, ?)`,
			KEYSPACE_TEST,
			ts.DBTable(),
		)
	default:
		return ""
	}
}

func cleanTable(ts TestSuite) {
	ts.T().Logf("cleaning test database after test '%s'", ts.T().Name())

	query := fmt.Sprintf(
		`TRUNCATE %s.%s`,
		KEYSPACE_TEST,
		ts.DBTable(),
	)

	err := ts.DBSession().Query(query).Exec()
	if err != nil {
		ts.FailNowf("unable to clean table", err.Error())
	}

	ts.T().Logf("test database cleaned after test '%s'", ts.T().Name())
}

func tearDownDatabase(ts TestSuite) {
	ts.T().Log("tearing down database")

	query := fmt.Sprintf(
		`DROP TABLE %s.%s`,
		KEYSPACE_TEST,
		ts.DBTable(),
	)

	err := ts.DBSession().Query(query).Exec()
	if err != nil {
		ts.FailNowf("unable to drop table", err.Error())
	}

	query = fmt.Sprintf(
		`DROP KEYSPACE %s`,
		KEYSPACE_TEST,
	)

	err = ts.DBSession().Query(query).Exec()
	if err != nil {
		ts.FailNowf("unable to drop test database keyspace", err.Error())
	}

	ts.DBSession().Close()
}
