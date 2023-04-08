package database

import "log"

func IncrementCourseClearCount(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.ratings SET completions=completions+1 WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}
