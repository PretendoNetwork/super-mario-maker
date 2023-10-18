package database

import "log"

func InsertBufferQueueData(dataID uint64, slot uint32, buffer []byte) {
	_, err := Postgres.Exec(`INSERT INTO pretendo_smm.buffer_queues( id, data_id, slot, buffer ) VALUES ( gen_random_uuid(), $1, $2, $3 ) ON CONFLICT DO NOTHING`, dataID, slot, buffer)
	if err != nil {
		log.Fatal(err)
	}
}
