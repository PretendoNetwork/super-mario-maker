package datastore_db

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
	"github.com/lib/pq"
)

// GetObjectRatingsWithSlotByDataIDs gets the ratings of many objects in a single query.
// Unlike GetObjectRatingsWithSlotByDataID this does NOT check if the objects are available,
// the caller is expected to have already filtered for that
func GetObjectRatingsWithSlotByDataIDs(dataIDs []types.UInt64) (map[types.UInt64]types.List[datastore_types.DataStoreRatingInfoWithSlot], *nex.Error) {
	ratings := make(map[types.UInt64]types.List[datastore_types.DataStoreRatingInfoWithSlot], len(dataIDs))

	if len(dataIDs) == 0 {
		return ratings, nil
	}

	ids := make([]int64, 0, len(dataIDs))
	for i := range dataIDs {
		ids = append(ids, int64(dataIDs[i]))
	}

	rows, err := database.Postgres.Query(`
		SELECT data_id, slot, total_value, count, initial_value
		FROM datastore.object_ratings
		WHERE data_id = ANY($1::bigint[])
		ORDER BY data_id, slot`,
		pq.Array(ids),
	)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var dataID types.UInt64

		rating := datastore_types.NewDataStoreRatingInfoWithSlot()
		rating.Rating = datastore_types.NewDataStoreRatingInfo()

		err := rows.Scan(&dataID, &rating.Slot, &rating.Rating.TotalValue, &rating.Rating.Count, &rating.Rating.InitialValue)
		if err != nil {
			globals.Logger.Error(err.Error())
			// TODO - Send more specific errors?
			return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
		}

		ratings[dataID] = append(ratings[dataID], rating)
	}

	if err := rows.Err(); err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return ratings, nil
}
