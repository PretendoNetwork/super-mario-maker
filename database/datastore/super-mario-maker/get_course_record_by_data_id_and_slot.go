package datastore_smm_db

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetCourseRecordByDataIDAndSlot(dataID uint64, slot uint8) (*datastore_super_mario_maker_types.DataStoreGetCourseRecordResult, uint32) {
	errCode := datastore_db.IsObjectAvailable(dataID)
	if errCode != 0 {
		return nil, errCode
	}

	courseRecord := datastore_super_mario_maker_types.NewDataStoreGetCourseRecordResult()
	courseRecord.DataID = dataID
	courseRecord.Slot = slot
	courseRecord.CreatedTime = nex.NewDateTime(0)
	courseRecord.UpdatedTime = nex.NewDateTime(0)

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
			return nil, nex.Errors.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nil, nex.Errors.DataStore.Unknown
	}

	courseRecord.CreatedTime.FromTimestamp(createdDate)
	courseRecord.UpdatedTime.FromTimestamp(updatedDate)

	return courseRecord, 0
}
