package cassandra

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func Init(keyspace string, replicaCount int, clusterAddresses ...string) (*gocql.Session, error) {
	var (
		err           error
		clusterConfig *gocql.ClusterConfig
	)

	clusterConfig, Session, err = connectToCassandra(clusterAddresses...)

	if err != nil {
		log.Printf("Failed to connect to Cassandra: %v", err)
		return nil, err
	}

	log.Println("Connected to Cassandra")

	err = createKeyspace(replicaCount, keyspace)
	if err != nil {
		log.Printf("Failed to create keyspace '%s' : %v", keyspace, err)
		return nil, err
	}

	// Reconnect to the keyspace
	clusterConfig.Keyspace = keyspace
	Session, err = clusterConfig.CreateSession()
	if err != nil {
		log.Printf("Failed to connect to keyspace '%s': %v", keyspace, err)
		return nil, err
	}

	log.Printf("Successfully connected to '%s' keyspace. Now you may want to run migrations!", keyspace)
	return Session, nil
}

func connectToCassandra(clusterAddresses ...string) (*gocql.ClusterConfig, *gocql.Session, error) {
	clusterConfig := gocql.NewCluster(clusterAddresses...)
	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, nil, err
	}

	return clusterConfig, session, nil
}

func createKeyspace(replicaCount int, keyspace string) error {
	createKeyspaceQuery := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': '%d'
		}`, keyspace, replicaCount)

	if err := Session.Query(createKeyspaceQuery).Exec(); err != nil {
		log.Printf("failed to create keyspace '%s': %v", keyspace, err)
		return err
	}
	fmt.Printf("Keyspace '%s' created\n", keyspace)
	return nil
}
