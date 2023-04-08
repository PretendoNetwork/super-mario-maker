package database

import "log"

func SetCoursePlayable(courseID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET playable=true WHERE data_id=?`, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}
