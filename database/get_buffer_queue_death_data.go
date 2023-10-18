package database

import (
	"database/sql"

	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetBufferQueueDeathData(dataID uint64) [][]byte {
	pBufferQueue := make([][]byte, 0)

	rows, err := Postgres.Query(`SELECT buffer FROM pretendo_smm.buffer_queues WHERE data_id=$1 AND slot=3`, dataID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return pBufferQueue
	}

	for rows.Next() {
		var buffer []byte

		err := rows.Scan(&buffer)
		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
			return pBufferQueue
		}

		pBufferQueue = append(pBufferQueue, buffer)
	}

	return pBufferQueue
}
