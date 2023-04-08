package database

import "log"

func CreateDataStoreIDGeneratorRow(nodeID int) {
	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_smm.generator_last_id(node_id, last_id) VALUES (?, ?) IF NOT EXISTS`, nodeID, 0).Exec(); err != nil {
		log.Fatal(err)
	}
}
