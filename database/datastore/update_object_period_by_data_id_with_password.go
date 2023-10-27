package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func UpdateObjectDataTypeByDataIDWithPassword(dataID uint64, period uint16, password uint64) uint32 {
	var updatePassword uint64
	var underReview bool

	err := database.Postgres.QueryRow(`SELECT update_password, under_review FROM datastore.objects WHERE data_id=$1 AND upload_completed=TRUE AND deleted=FALSE`, dataID).Scan(
		&updatePassword,
		&underReview,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nex.Errors.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nex.Errors.DataStore.Unknown
	}

	if updatePassword != 0 && updatePassword != password {
		return nex.Errors.DataStore.InvalidPassword
	}

	if underReview {
		return nex.Errors.DataStore.UnderReviewing
	}

	_, err = database.Postgres.Exec(`UPDATE datastore.objects SET period=$1 WHERE data_id=$2`, period, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.Errors.DataStore.Unknown
	}

	return 0
}
