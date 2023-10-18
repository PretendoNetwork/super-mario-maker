package database

import "log"

func IncrementCourseAttemptCount(courseID uint64) {
	_, err := Postgres.Exec(`UPDATE pretendo_smm.ratings SET attempts=attempts+1 WHERE data_id=$1`, courseID)
	if err != nil {
		log.Fatal(err)
	}
}
