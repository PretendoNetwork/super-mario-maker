package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetObjectRatingsWithSlotByDataID(dataID uint64) ([]*datastore_types.DataStoreRatingInfoWithSlot, uint32) {
	errCode := IsObjectAvailable(dataID)
	if errCode != 0 {
		return nil, errCode
	}

	ratings := make([]*datastore_types.DataStoreRatingInfoWithSlot, 0)

	rows, err := database.Postgres.Query(`SELECT slot, total_value, count, initial_value FROM datastore.object_ratings WHERE data_id=$1`, dataID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nex.Errors.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nil, nex.Errors.DataStore.Unknown
	}

	defer rows.Close()

	for rows.Next() {
		rating := datastore_types.NewDataStoreRatingInfoWithSlot()
		rating.Rating = datastore_types.NewDataStoreRatingInfo()

		err := rows.Scan(&rating.Slot, &rating.Rating.TotalValue, &rating.Rating.Count, &rating.Rating.InitialValue)

		if err != nil {
			globals.Logger.Error(err.Error())
			// TODO - Send more specific errors?
			return nil, nex.Errors.DataStore.Unknown
		}

		ratings = append(ratings, rating)
	}

	if err := rows.Err(); err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.Errors.DataStore.Unknown
	}

	return ratings, 0
}
