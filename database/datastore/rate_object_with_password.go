package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func RateObjectWithPassword(dataID types.UInt64, slot types.UInt8, ratingValue types.Int32, accessPassword types.UInt64) (datastore_types.DataStoreRatingInfo, *nex.Error) {
	nexError := IsObjectAvailableWithPassword(dataID, accessPassword)
	if nexError != nil {
		return datastore_types.NewDataStoreRatingInfo(), nexError
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
			return datastore_types.NewDataStoreRatingInfo(), nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, "Invalid argument")
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return datastore_types.NewDataStoreRatingInfo(), nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return rating, nil
}
