package database

import "log"

func IncrementCourseFailCount(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET failures=failures+1 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}
