package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func UpdateObjectUploadCompletedByDataID(dataID uint64, uploadCompleted bool) uint32 {
	var underReview bool

	err := database.Postgres.QueryRow(`SELECT update_password FROM datastore.objects WHERE data_id=$1 AND deleted=FALSE`, dataID).Scan(&underReview)
	if err != nil {
		if err == sql.ErrNoRows {
			return nex.Errors.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nex.Errors.DataStore.Unknown
	}

	if underReview {
		return nex.Errors.DataStore.UnderReviewing
	}

	_, err = database.Postgres.Exec(`UPDATE datastore.objects SET upload_completed=$1 WHERE data_id=$2`, uploadCompleted, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.Errors.DataStore.Unknown
	}

	return 0
}
