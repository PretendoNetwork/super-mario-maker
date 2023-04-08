package database

func GetDataStoreIDGeneratorLastID(nodeID int) uint32 {
	var lastID uint32
	_ = cassandraClusterSession.Query(`SELECT last_id FROM pretendo_smm.generator_last_id WHERE node_id=?`, nodeID).Scan(&lastID)

	return lastID
}
