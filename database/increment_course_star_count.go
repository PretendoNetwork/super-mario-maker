package database

import "log"

func IncrementCourseStarCount(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET stars=stars+1 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}
