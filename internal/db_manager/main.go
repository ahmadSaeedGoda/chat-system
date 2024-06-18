// The purpose of this package is to handle several different DB connections.
// e.g PostgreSQL for configs, Cassandra for app data, ES for app usage behaviors and Data-Lakes, etc
// Many different DB engines can be employed for diff purposes, where each one is best suited for some data nature

package dbmanager

import (
	"chat-system/internal/cassandra"
	"log"
	"os"
	"strings"

	"github.com/gocql/gocql"
)

const CASSANDRA_KEYSPACE = "chat"

var CassandraSession *gocql.Session

func InitDB() *gocql.Session {
	var err error
	CassandraSession, err = setupCassandra()
	if err != nil {
		log.Fatalf("error while setting up cassandra: %v", err)
	}

	return CassandraSession
}

func setupCassandra() (*gocql.Session, error) {
	cassandraNodes := os.Getenv("CASSANDRA_NODES")
	clusterAddresses := strings.Split(cassandraNodes, ",")
	replicaCount := len(clusterAddresses)

	return cassandra.Init(CASSANDRA_KEYSPACE, replicaCount, clusterAddresses...)
}
