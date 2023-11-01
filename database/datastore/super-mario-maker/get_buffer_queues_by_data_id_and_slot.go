package datastore_smm_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
)

func GetBufferQueuesByDataIDAndSlot(dataID uint64, slot uint32) ([][]byte, uint32) {
	errCode := datastore_db.IsObjectAvailable(dataID)
	if errCode != 0 {
		return nil, errCode
	}

	bufferQueues := make([][]byte, 0)

	rows, err := database.Postgres.Query(`SELECT buffer FROM datastore.buffer_queues WHERE data_id=$1 AND slot=$2 ORDER BY creation_date`, dataID, slot)

	// * No rows is allowed
	if err != nil && err != sql.ErrNoRows {
		// TODO - Send more specific errors?
		return nil, nex.Errors.DataStore.Unknown
	}

	defer rows.Close()

	for rows.Next() {
		var buffer []byte

		err := rows.Scan(&buffer)
		if err != nil {
			// TODO - Send more specific errors?
			return nil, nex.Errors.DataStore.Unknown
		}

		bufferQueues = append(bufferQueues, buffer)
	}

	if err := rows.Err(); err != nil {
		// TODO - Send more specific errors?
		return nil, nex.Errors.DataStore.Unknown
	}

	return bufferQueues, 0
}
