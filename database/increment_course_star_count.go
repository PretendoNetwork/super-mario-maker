package database

import "log"

func IncrementCourseStarCount(courseID uint64) {
	_, err := Postgres.Exec(`UPDATE pretendo_smm.ratings SET stars=stars+1 WHERE data_id=$1`, courseID)
	if err != nil {
		log.Fatal(err)
	}
}
