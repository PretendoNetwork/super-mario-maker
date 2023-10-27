package datastore_smm_db

import (
	"database/sql"

	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/lib/pq"
)

func GetCustomRankingsByDataIDs(applicationID uint32, dataIDs []uint64) []*datastore_super_mario_maker_types.DataStoreCustomRankingResult {
	// * Unlike other methods which query for multiple objects,
	// * DataStoreSMM::GetCustomRankingByDataID *OMITS* objects
	// * from it's response if they could not be found, rather
	// * than using a zero-ed object and DataStore::NotFound as
	// * the Result. Because of this, all errors are just thrown
	// * away here and not sent back to the client
	results := make([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, 0, len(dataIDs))

	rows, err := database.Postgres.Query(`SELECT
		data_id,
		value
	FROM datastore.object_custom_rankings WHERE application_id=$1 AND data_id=ANY($2)`,
		applicationID,
		pq.Array(dataIDs),
	)

	// * No rows is allowed
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Error(err.Error())
		return results
	}

	defer rows.Close()

	for rows.Next() {
		var dataID uint64
		var value uint32

		err := rows.Scan(&dataID, &value)
		if err != nil {
			globals.Logger.Error(err.Error())
			continue
		}

		objectInfo, errCode := datastore_db.GetObjectInfoByDataID(dataID)
		if errCode != 0 {
			globals.Logger.Errorf("Got error code %d for object %d", errCode, dataID)
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
