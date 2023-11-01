package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func RateObjectWithPassword(dataID uint64, slot uint8, ratingValue int32, accessPassword uint64) (*datastore_types.DataStoreRatingInfo, uint32) {
	errCode := IsObjectAvailableWithPassword(dataID, accessPassword)
	if errCode != 0 {
		return nil, errCode
	}

	rating := datastore_types.NewDataStoreRatingInfo()

	err := database.Postgres.QueryRow(`
		UPDATE datastore.object_ratings
		SET total_value=total_value+$1, count=count+1
		WHERE data_id=$2 AND slot=$3
		RETURNING total_value, count, initial_value`, ratingValue, dataID, slot,
	).Scan(
		&rating.TotalValue,
		&rating.Count,
		&rating.InitialValue,
	)

	if err != nil {
		// * If no data was returned, that means the slot
		// * was out of bounds since the object exists due
		// * to the check at the start of the function.
		// * This is an invalid argument
		if err == sql.ErrNoRows {
			return nil, nex.Errors.DataStore.InvalidArgument
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.Errors.DataStore.Unknown
	}

	return rating, 0
}
