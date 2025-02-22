package datastore_smm_db

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func GetCourseRecordByDataIDAndSlot(dataID types.UInt64, slot types.UInt8) (datastore_super_mario_maker_types.DataStoreGetCourseRecordResult, *nex.Error) {
	nexError := datastore_db.IsObjectAvailable(dataID)
	if nexError != nil {
		return datastore_super_mario_maker_types.NewDataStoreGetCourseRecordResult(), nexError
	}

	courseRecord := datastore_super_mario_maker_types.NewDataStoreGetCourseRecordResult()
	courseRecord.DataID = dataID
	courseRecord.Slot = slot
	var createdDate time.Time
	var updatedDate time.Time

	err := database.Postgres.QueryRow(`SELECT
		first_pid,
		best_pid,
		best_score,
		creation_date,
		update_date
	FROM datastore.course_records WHERE data_id=$1 AND slot=$2`,
		dataID,
		slot,
	).Scan(
		&courseRecord.FirstPID,
		&courseRecord.BestPID,
		&courseRecord.BestScore,
		&createdDate,
		&updatedDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return datastore_super_mario_maker_types.NewDataStoreGetCourseRecordResult(), nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return datastore_super_mario_maker_types.NewDataStoreGetCourseRecordResult(), nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	courseRecord.CreatedTime.FromTimestamp(createdDate)
	courseRecord.UpdatedTime.FromTimestamp(updatedDate)

	return courseRecord, nil
}
