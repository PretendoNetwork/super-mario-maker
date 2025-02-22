package datastore_smm_db

import (
	"log"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
)

func InsertOrUpdateBufferQueueData(dataID types.UInt64, slot types.UInt32, buffer types.QBuffer) *nex.Error {
	nexError := datastore_db.IsObjectAvailable(dataID)
	if nexError != nil {
		return nexError
	}

	now := time.Now()

	// * Real server does not allow duplicate
	// * buffers to be in a given objects slot,
	// * even if the buffers were uploaded by
	// * different clients. Instead, it removes
	// * the older buffer and adds the newer one.
	// * Instead of that, we just update the
	// * creation time
	_, err := database.Postgres.Exec(`INSERT INTO datastore.buffer_queues (
		data_id,
		slot,
		creation_date,
		buffer
	) VALUES (
		$1,
		$2,
		$3,
		$4
	) ON CONFLICT (data_id, slot, buffer) DO UPDATE SET creation_date=$3`, dataID, slot, now, buffer)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
