package datastore_smm_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker/globals"
	"github.com/lib/pq"
)

func GetCustomRankingsByDataIDs(applicationID types.UInt32, dataIDs types.List[types.UInt64]) types.List[datastore_super_mario_maker_types.DataStoreCustomRankingResult] {
	// * Unlike other methods which query for multiple objects,
	// * DataStoreSMM::GetCustomRankingByDataID *OMITS* objects
	// * from it's response if they could not be found, rather
	// * than using a zero-ed object and DataStore::NotFound as
	// * the Result. Because of this, all errors are just thrown
	// * away here and not sent back to the client
	results := make(types.List[datastore_super_mario_maker_types.DataStoreCustomRankingResult], 0, len(dataIDs))

	// * Using UNNEST and WITH ORDINALITY because the input
	// * array may contain duplicate DataIDs. These duplicate
	// * DataIDs should result in duplicate results, and this
	// * was the most efficient way to do this as it doesn't
	// * involve processing all rows or manually duplicating
	// * rows
	rows, err := database.Postgres.Query(`
	SELECT
		rankings.data_id,
		rankings.value
	FROM datastore.object_custom_rankings rankings
	JOIN UNNEST($1::bigint[])
	WITH ORDINALITY AS rows(data_id, ord)
		ON rankings.data_id = rows.data_id
		AND rankings.application_id = $2
	ORDER BY rows.ord`,
		pq.Array(dataIDs),
		applicationID,
	)

	// * No rows is allowed
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Error(err.Error())
		return results
	}

	defer rows.Close()

	for rows.Next() {
		var dataID types.UInt64
		var value types.UInt32

		err := rows.Scan(&dataID, &value)
		if err != nil {
			globals.Logger.Error(err.Error())
			continue
		}

		objectInfo, nexError := datastore_db.GetObjectInfoByDataID(dataID)
		if nexError != nil {
			globals.Logger.Errorf("Got error code %d for object %d", nexError.ResultCode, dataID)
			continue
		}

		result := datastore_super_mario_maker_types.NewDataStoreCustomRankingResult()

		// * Order is always 0, for some reason
		result.Score = value
		result.MetaInfo = objectInfo

		results = append(results, result)
	}

	return results
}
