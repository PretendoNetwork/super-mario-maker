package database

import (
	"log"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/types"
)

func GetCourseMetadatasByLimit(limit uint32) []*types.CourseMetadata {
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT data_id, owner_pid, size, name, meta_binary, flag, creation_date, update_date, data_type, period FROM pretendo_smm.courses LIMIT ?`, limit).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	courseMetadatas := make([]*types.CourseMetadata, 0)

	for _, course := range sliceMap {
		dataID := uint64(course["data_id"].(int64))

		var stars uint32
		var attempts uint32
		var failures uint32
		var completions uint32

		_ = cassandraClusterSession.Query(`SELECT stars, attempts, failures, completions FROM pretendo_smm.ratings WHERE data_id=?`, dataID).Scan(&stars, &attempts, &failures, &completions)

		courseMetadata := &types.CourseMetadata{
			DataID:      dataID,
			OwnerPID:    uint32(course["owner_pid"].(int)),
			Size:        uint32(course["size"].(int)),
			CreatedTime: nex.NewDateTime(uint64(course["creation_date"].(int64))),
			UpdatedTime: nex.NewDateTime(uint64(course["update_date"].(int64))),
			Name:        course["name"].(string),
			MetaBinary:  course["meta_binary"].([]byte),
			Stars:       stars,
			Attempts:    attempts,
			Failures:    failures,
			Completions: completions,
			Flag:        uint32(course["flag"].(int)),
			DataType:    uint16(course["data_type"].(int16)),
			Period:      uint16(course["period"].(int16)),
		}

		courseMetadatas = append(courseMetadatas, courseMetadata)
	}

	return courseMetadatas
}
