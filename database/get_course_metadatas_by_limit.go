package database

import (
	"database/sql"
	"log"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/types"
)

func GetCourseMetadatasByLimit(limit uint32) []*types.CourseMetadata {
	rows, err := Postgres.Query(`SELECT data_id, owner_pid, size, name, meta_binary, flag, creation_date, update_date, data_type, period FROM pretendo_smm.courses LIMIT $1`, limit)
	if err != nil {
		log.Fatal(err)
	}

	courseMetadatas := make([]*types.CourseMetadata, 0)

	for rows.Next() {
		var dataID uint64
		var ownerPID uint32
		var size uint32
		var name string
		var metaBinary []byte
		var flag uint32
		var creationDate uint64
		var updateDate uint64
		var dataType uint16
		var period uint16

		err := rows.Scan(
			&dataID,
			&ownerPID,
			&size,
			&name,
			&metaBinary,
			&flag,
			&creationDate,
			&updateDate,
			&dataType,
			&period,
		)
		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
		}

		var stars uint32
		var attempts uint32
		var failures uint32
		var completions uint32

		_ = Postgres.QueryRow(`SELECT stars, attempts, failures, completions FROM pretendo_smm.ratings WHERE data_id=$1`, dataID).Scan(&stars, &attempts, &failures, &completions)

		courseMetadata := &types.CourseMetadata{
			DataID:      dataID,
			OwnerPID:    ownerPID,
			Size:        size,
			CreatedTime: nex.NewDateTime(creationDate),
			UpdatedTime: nex.NewDateTime(updateDate),
			Name:        name,
			MetaBinary:  metaBinary,
			Stars:       stars,
			Attempts:    attempts,
			Failures:    failures,
			Completions: completions,
			Flag:        flag,
			DataType:    dataType,
			Period:      period,
		}

		courseMetadatas = append(courseMetadatas, courseMetadata)
	}

	return courseMetadatas
}
