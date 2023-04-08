package database

import (
	"log"

	"github.com/PretendoNetwork/nex-go"
)

func UpdateCourseWorldRecord(courseID uint64, ownerPID uint32, score int32) {
	datetime := nex.NewDateTime(0)
	now := datetime.Now()

	if GetCourseWorldRecord(courseID) == nil {
		if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET world_record_first_pid=?, world_record_creation_date=? WHERE data_id=?`, ownerPID, now, courseID).Exec(); err != nil {
			log.Fatal(err)
		}
	}

	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET world_record_pid=?, world_record_update_date=?, world_record=? WHERE data_id=?`, ownerPID, now, score, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}
