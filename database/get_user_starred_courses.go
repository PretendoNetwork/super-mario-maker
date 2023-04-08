package database

import "github.com/PretendoNetwork/super-mario-maker-secure/types"

func GetUserStarredCourses(pid uint32) []*types.CourseMetadata {
	var dataIDs []uint64
	_ = cassandraClusterSession.Query(`SELECT starred_courses FROM pretendo_smm.user_play_info WHERE pid=?`, pid).Scan(&dataIDs)

	return GetCourseMetadataByDataIDs(dataIDs)
}
