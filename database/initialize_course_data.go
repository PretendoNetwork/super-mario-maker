package database

import (
	"log"

	"github.com/PretendoNetwork/nex-go"
	"github.com/lib/pq"
)

func InitializeCourseData(ownerPID uint32, size uint32, name string, flag uint32, extraData []string, dataType uint16, period uint16) uint64 {
	datetime := nex.NewDateTime(0)
	now := datetime.Now()
	var dataID uint64

	err := Postgres.QueryRow(`INSERT INTO pretendo_smm.courses(
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
		$1,
		$2,
		$3,
		$4,
		$5,
		false,
		$6,
		$7,
		0,
		0,
		0,
		0,
		0,
		$8,
		$9
	) RETURNING data_id`, ownerPID, size, name, flag, pq.Array(extraData), now, now, dataType, period).Scan(&dataID)

	if err != nil {
		log.Fatal(err)
	}

	_, err = Postgres.Exec(`INSERT INTO pretendo_smm.ratings (
		data_id,
		stars,
		attempts,
		failures,
		completions
	)
	VALUES (
		$1,
		0,
		0,
		0,
		0
	)`, dataID)
	if err != nil {
		log.Fatal(err)
	}

	return dataID
}
