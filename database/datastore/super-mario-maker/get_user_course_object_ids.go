package datastore_smm_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetUserCourseObjectIDs(ownerPID uint32) ([]uint64, uint32) {
	courseObjectIDs := make([]uint64, 0)

	// * Course objects seem to have data types > 2 and < 50.
	// * Data type 1 seems to be reserved for "maker" objects.
	// * Data type 2 seems to be reserved for objects
	// * created through "PrepareAttachFile".
	// * Data type 50 is reserved for the Event Courses metadata
	// * file, and data type 51 is reserved for event courses
	rows, err := database.Postgres.Query(`SELECT data_id FROM datastore.objects WHERE owner=$1 AND data_type > 2 AND data_type < 50`, ownerPID)

	// * No rows is allowed
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.DataStore.Unknown
	}

	defer rows.Close()

	for rows.Next() {
		var dataID uint64

		err := rows.Scan(&dataID)
		if err != nil {
			globals.Logger.Error(err.Error())
			continue
		}

		errCode := datastore_db.IsObjectAvailable(dataID)
		if errCode != 0 {
			continue
		}

		courseObjectIDs = append(courseObjectIDs, dataID)
	}

	if err := rows.Err(); err != nil {
		// TODO - Send more specific errors?
		return nil, nex.Errors.DataStore.Unknown
	}

	return courseObjectIDs, 0
}
