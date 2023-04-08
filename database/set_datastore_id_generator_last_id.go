package database

import "log"

func SetDataStoreIDGeneratorLastID(nodeID int, value uint32) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.generator_last_id SET last_id=? WHERE node_id=?`, value, nodeID).Exec(); err != nil {
		log.Fatal(err)
	}
}
