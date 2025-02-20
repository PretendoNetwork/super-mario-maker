package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func UpdateObjectMetaBinaryByDataIDWithPassword(dataID types.UInt64, metaBinary types.QBuffer, password types.UInt64) *nex.Error {
	var updatePassword types.UInt64
	var underReview bool

	err := database.Postgres.QueryRow(`SELECT update_password, under_review FROM datastore.objects WHERE data_id=$1 AND upload_completed=TRUE AND deleted=FALSE`, dataID).Scan(
		&updatePassword,
		&underReview,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	if updatePassword != 0 && updatePassword != password {
		return nex.NewError(nex.ResultCodes.DataStore.InvalidPassword, "Invalid password")
	}

	if underReview {
		return nex.NewError(nex.ResultCodes.DataStore.UnderReviewing, "This object is under review")
	}

	_, err = database.Postgres.Exec(`UPDATE datastore.objects SET meta_binary=$1 WHERE data_id=$2`, metaBinary, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return nil
}
