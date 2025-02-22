package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func IsObjectAvailable(dataID types.UInt64) *nex.Error {
	var underReview bool

	err := database.Postgres.QueryRow(`SELECT under_review FROM datastore.objects WHERE data_id=$1 AND upload_completed=TRUE AND deleted=FALSE`, dataID).Scan(&underReview)
	if err != nil {
		if err == sql.ErrNoRows {
			return nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	if underReview {
		return nex.NewError(nex.ResultCodes.DataStore.UnderReviewing, "This object is currently under review")
	}

	return nil
}
