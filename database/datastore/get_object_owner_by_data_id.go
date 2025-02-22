package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func GetObjectOwnerByDataID(dataID types.UInt64) (uint32, *nex.Error) {
	var owner uint32

	err := database.Postgres.QueryRow(`SELECT owner FROM datastore.objects WHERE data_id=$1`, dataID).Scan(&owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return 0, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return owner, nil
}
