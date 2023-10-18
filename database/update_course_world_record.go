package database

import (
	"log"

	"github.com/PretendoNetwork/nex-go"
)

func UpdateCourseWorldRecord(courseID uint64, ownerPID uint32, score int32) {
	datetime := nex.NewDateTime(0)
	now := datetime.Now()

	if GetCourseWorldRecord(courseID) == nil {
		_, err := Postgres.Exec(`UPDATE pretendo_smm.courses SET world_record_first_pid=$1, world_record_creation_date=$2 WHERE data_id=$3`, ownerPID, now, courseID)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err := Postgres.Exec(`UPDATE pretendo_smm.courses SET world_record_pid=$1, world_record_update_date=$2, world_record=$3 WHERE data_id=$4`, ownerPID, now, score, courseID)
	if err != nil {
		log.Fatal(err)
	}
}
