package database

import "log"

func GetBufferQueueDeathData(dataID uint64) [][]byte {
	pBufferQueue := make([][]byte, 0)

	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT buffer FROM pretendo_smm.buffer_queues WHERE data_id=? AND slot=3 ALLOW FILTERING`, dataID).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	for _, bufferQueue := range sliceMap {
		pBufferQueue = append(pBufferQueue, bufferQueue["buffer"].([]byte))
	}

	return pBufferQueue
}
