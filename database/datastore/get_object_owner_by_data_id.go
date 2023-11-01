package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetObjectOwnerByDataID(dataID uint64) (uint32, uint32) {
	var owner uint32

	err := database.Postgres.QueryRow(`SELECT owner FROM datastore.objects WHERE data_id=$1`, dataID).Scan(&owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nex.Errors.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return 0, nex.Errors.DataStore.Unknown
	}

	return owner, 0
}
