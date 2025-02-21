package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func GetObjectRatingsWithSlotByDataID(dataID types.UInt64) ([]datastore_types.DataStoreRatingInfoWithSlot, *nex.Error) {
	nexError := IsObjectAvailable(dataID)
	if nexError != nil {
		return nil, nexError
	}

	ratings := types.NewList[datastore_types.DataStoreRatingInfoWithSlot]()

	rows, err := database.Postgres.Query(`SELECT slot, total_value, count, initial_value FROM datastore.object_ratings WHERE data_id=$1`, dataID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		rating := datastore_types.NewDataStoreRatingInfoWithSlot()
		rating.Rating = datastore_types.NewDataStoreRatingInfo()

		err := rows.Scan(&rating.Slot, &rating.Rating.TotalValue, &rating.Rating.Count, &rating.Rating.InitialValue)

		if err != nil {
			globals.Logger.Error(err.Error())
			// TODO - Send more specific errors?
			return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
		}

		ratings = append(ratings, rating)
	}

	if err := rows.Err(); err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return ratings, nil
}
