package database

import (
	"database/sql"
	"log"

	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/types"
)

func GetCourseMetadatasByPID(pid uint32) []*types.CourseMetadata {
	courseMetadatas := make([]*types.CourseMetadata, 0)

	rows, err := Postgres.Query(`SELECT data_id FROM pretendo_smm.courses WHERE owner_pid=$1`, pid)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var dataID uint64

		err := rows.Scan(&dataID)
		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
			return courseMetadatas
		}

		courseMetadatas = append(courseMetadatas, GetCourseMetadataByDataID(dataID))
	}

	return courseMetadatas
}
