package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocql/gocql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const KEYSPACE = "chat"

var Session *gocql.Session

func Init() {
	cassandraNodes := os.Getenv("CASSANDRA_NODES")
	clusterAddresses := strings.Split(cassandraNodes, ",")

	var (
		err           error
		clusterConfig *gocql.ClusterConfig
	)

	clusterConfig, Session, err = connectToCassandra(clusterAddresses...)

	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	log.Println("Connected to Cassandra")

	createKeyspace(len(clusterAddresses))

	// Reconnect to the keyspace
	clusterConfig.Keyspace = "chat"
	Session, err = clusterConfig.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to keyspace: %v", err)
	}

	log.Printf("Successfully connected to '%s' keyspace. Now you may want to run migrations!", KEYSPACE)
}

func connectToCassandra(clusterAddresses ...string) (*gocql.ClusterConfig, *gocql.Session, error) {
	clusterConfig := gocql.NewCluster(clusterAddresses...)
	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, nil, err
	}

	return clusterConfig, session, nil
}

func createKeyspace(replicaCount int) {
	createKeyspaceQuery := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': '%d'
		}`, KEYSPACE, replicaCount)

	if err := Session.Query(createKeyspaceQuery).Exec(); err != nil {
		log.Fatalf("failed to create keyspace: %v", err)
	}
	fmt.Printf("Keyspace '%s' created\n", KEYSPACE)
}
