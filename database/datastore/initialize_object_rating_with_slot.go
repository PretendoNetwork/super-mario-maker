package datastore_db

import (
	"github.com/PretendoNetwork/nex-go/v2"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func InitializeObjectRatingWithSlot(dataID uint64, param datastore_types.DataStoreRatingInitParamWithSlot) *nex.Error {
	_, err := database.Postgres.Exec(`INSERT INTO datastore.object_ratings (
		data_id,
		slot,
		flag,
		internal_flag,
		lock_type,
		initial_value,
		range_min,
		range_max,
		period_hour,
		period_duration,
		total_value
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11
	)`,
		dataID,
		param.Slot,
		param.Param.Flag,
		param.Param.InternalFlag,
		param.Param.LockType,
		param.Param.InitialValue,
		param.Param.RangeMin,
		param.Param.RangeMax,
		param.Param.PeriodHour,
		param.Param.PeriodDuration,
		param.Param.InitialValue, // * Start the value off at the initial value
	)

	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return nil
}
