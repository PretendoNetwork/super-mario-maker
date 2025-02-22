package datastore_smm_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
)

func GetBufferQueuesByDataIDAndSlot(dataID types.UInt64, slot types.UInt32) (types.List[types.QBuffer], *nex.Error) {
	nexError := datastore_db.IsObjectAvailable(dataID)
	if nexError != nil {
		return nil, nexError
	}

	bufferQueues := types.NewList[types.QBuffer]()

	rows, err := database.Postgres.Query(`SELECT buffer FROM datastore.buffer_queues WHERE data_id=$1 AND slot=$2 ORDER BY creation_date`, dataID, slot)

	// * No rows is allowed
	if err != nil && err != sql.ErrNoRows {
		// TODO - Send more specific errors?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var buffer types.QBuffer

		err := rows.Scan(&buffer)
		if err != nil {
			// TODO - Send more specific errors?
			return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
		}

		bufferQueues = append(bufferQueues, buffer)
	}

	if err := rows.Err(); err != nil {
		// TODO - Send more specific errors?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return bufferQueues, nil
}
