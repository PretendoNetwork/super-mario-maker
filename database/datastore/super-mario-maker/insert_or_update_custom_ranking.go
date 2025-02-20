package datastore_smm_db

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func InsertOrUpdateCustomRanking(dataID types.UInt64, applicationID, score types.UInt32) *nex.Error {
	nexError := datastore_db.IsObjectAvailable(dataID)
	if nexError != nil {
		globals.Logger.Errorf("Error code %d", nexError.ResultCode)
		return nexError
	}

	_, err := database.Postgres.Exec(`INSERT INTO datastore.object_custom_rankings (
		data_id,
		application_id,
		value
	) VALUES (
		$1,
		$2,
		$3
	) ON CONFLICT (data_id, application_id) DO UPDATE SET value=datastore.object_custom_rankings.value+EXCLUDED.value`,
		dataID,
		applicationID,
		score,
	)

	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return nil
}
