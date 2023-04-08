package database

import (
	"log"

	"github.com/PretendoNetwork/super-mario-maker-secure/types"
)

func GetCourseMetadatasByPID(pid uint32) []*types.CourseMetadata {
	courseMetadatas := make([]*types.CourseMetadata, 0)

	// TODO: Fix this query? Seems like a weird way of doing this...
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT data_id FROM pretendo_smm.courses WHERE owner_pid=? ALLOW FILTERING`, pid).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	for _, course := range sliceMap {
		dataID := uint64(course["data_id"].(int64))
		courseMetadatas = append(courseMetadatas, GetCourseMetadataByDataID(dataID))
	}

	return courseMetadatas
}
