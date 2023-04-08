package database

import "log"

func InsertBufferQueueData(dataID uint64, slot uint32, buffer []byte) {
	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_smm.buffer_queues( id, data_id, slot, buffer ) VALUES ( now(), ?, ?, ? ) IF NOT EXISTS`, dataID, slot, buffer).Exec(); err != nil {
		log.Fatal(err)
	}
}
