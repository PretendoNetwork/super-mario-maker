package datastore_smm_db

import (
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func InsertOrUpdateCourseRecord(dataID types.UInt64, slot types.UInt8, pid types.PID, score types.Int32) *nex.Error {
	nexError := datastore_db.IsObjectAvailable(dataID)
	if nexError != nil {
		globals.Logger.Errorf("Error code %d", nexError.ResultCode)
		return nexError
	}

	now := time.Now()

	_, err := database.Postgres.Exec(`INSERT INTO datastore.course_records (
		data_id,
		slot,
		first_pid,
		best_pid,
		best_score,
		creation_date,
		update_date
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	) ON CONFLICT (data_id, slot) DO UPDATE
	SET best_score = CASE WHEN datastore.course_records.best_score > $5 THEN $5 ELSE datastore.course_records.best_score END,
		best_pid = CASE WHEN datastore.course_records.best_score > $5 THEN $4 ELSE datastore.course_records.best_pid END,
		update_date = CASE WHEN datastore.course_records.best_score > $5 THEN $7 ELSE datastore.course_records.update_date END`,
		dataID,
		slot,
		pid,
		pid,
		score,
		now,
		now,
	)

	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return nil
}
