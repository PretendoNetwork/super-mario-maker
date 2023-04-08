package database

import "log"

func UpdateCourseMetaBinary(courseID uint64, metaBinary []byte) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_smm.courses SET meta_binary=? WHERE data_id=?`, metaBinary, courseID).Exec(); err != nil {
		log.Fatal(err)
	}
}
