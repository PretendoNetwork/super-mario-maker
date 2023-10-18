package database

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/types"
)

func GetCourseWorldRecord(dataID uint64) *types.CourseWorldRecord {
	var worldRecordFirstPID uint32
	var worldRecordPID uint32
	var worldRecordCreatedTime uint64
	var worldRecordUpdatedTime uint64
	var worldRecord int32

	_ = Postgres.QueryRow(`SELECT world_record_first_pid, world_record_pid, world_record_creation_date, world_record_update_date, world_record FROM pretendo_smm.courses WHERE data_id=$1`, dataID).Scan(&worldRecordFirstPID, &worldRecordPID, &worldRecordCreatedTime, &worldRecordUpdatedTime, &worldRecord)

	if worldRecordFirstPID == 0 {
		return nil
	}

	return &types.CourseWorldRecord{
		FirstPID:    worldRecordFirstPID,
		BestPID:     worldRecordPID,
		CreatedTime: nex.NewDateTime(worldRecordCreatedTime),
		UpdatedTime: nex.NewDateTime(worldRecordUpdatedTime),
		Score:       worldRecord,
	}
}
