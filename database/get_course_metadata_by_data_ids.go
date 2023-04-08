package database

import "github.com/PretendoNetwork/super-mario-maker-secure/types"

func GetCourseMetadataByDataIDs(dataIDs []uint64) []*types.CourseMetadata {
	// TODO: Do this in one query?
	courseMetadatas := make([]*types.CourseMetadata, 0)

	for _, dataID := range dataIDs {
		courseMetadata := GetCourseMetadataByDataID(dataID)

		if courseMetadata != nil {
			courseMetadatas = append(courseMetadatas, courseMetadata)
		}
	}

	return courseMetadatas
}
