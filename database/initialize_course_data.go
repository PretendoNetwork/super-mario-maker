package database

import (
	"log"

	"github.com/PretendoNetwork/nex-go"
)

func InitializeCourseData(courseID uint64, ownerPID uint32, size uint32, name string, flag uint32, extraData []string, dataType uint16, period uint16) {
	datetime := nex.NewDateTime(0)
	now := datetime.Now()

	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_smm.courses(
		data_id,
		owner_pid,
		size,
		name,
		flag,
		extra_data,
		playable,
		creation_date,
		update_date,
		world_record_first_pid,
		world_record_pid,
		world_record_creation_date,
		world_record_update_date,
		world_record,
		data_type,
		period
	)
	VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		false,
		?,
		?,
		0,
		0,
		0,
		0,
		0,
		?,
		?
	) IF NOT EXISTS`,
		courseID, ownerPID, size, name, flag, extraData, now, now, dataType, period).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET stars=stars+0, attempts=attempts+0, failures=failures+0, completions=completions+0 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}
