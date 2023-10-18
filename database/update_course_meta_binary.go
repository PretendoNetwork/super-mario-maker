package database

import "log"

func UpdateCourseMetaBinary(courseID uint64, metaBinary []byte) {
	_, err := Postgres.Exec(`UPDATE pretendo_smm.courses SET meta_binary=$1 WHERE data_id=$2`, metaBinary, courseID)
	if err != nil {
		log.Fatal(err)
	}
}
