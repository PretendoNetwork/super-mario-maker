package database

import "log"

func SetCoursePlayable(courseID uint64) {
	_, err := Postgres.Exec(`UPDATE pretendo_smm.courses SET playable=true WHERE data_id=$1`, courseID)
	if err != nil {
		log.Fatal(err)
	}
}
