package database

import (
	"log"

	"github.com/PretendoNetwork/super-mario-maker-secure/types"
)

func GetUserStarredCourses(pid uint32) []*types.CourseMetadata {
	var dataIDs []uint64

	rows, err := Postgres.Query(`SELECT starred_courses FROM pretendo_smm.user_play_info WHERE pid=$1`, pid)
	if err != nil {
		log.Fatal(err)
	}

	rows.Scan(&dataIDs)

	return GetCourseMetadataByDataIDs(dataIDs)
}
