package database

import "log"

func IncrementCourseAttemptCount(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET attempts=attempts+1 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}
