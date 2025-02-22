package datastore_db

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func DeleteObjectByDataID(dataID types.UInt64) *nex.Error {
	nexError := IsObjectAvailable(dataID)
	if nexError != nil {
		return nexError
	}

	_, err := database.Postgres.Exec(`UPDATE datastore.objects SET deleted=TRUE WHERE data_id=$1`, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return nil
}
