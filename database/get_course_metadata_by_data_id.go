package database

import (
	"log"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/types"
	"github.com/gocql/gocql"
)

func GetCourseMetadataByDataID(dataID uint64) *types.CourseMetadata {
	var ownerPID uint32
	var size uint32
	var name string
	var metaBinary []byte
	var flag uint32
	var createdTime uint64
	var updatedTime uint64
	var dataType uint16
	var period uint16

	err := cassandraClusterSession.Query(`SELECT owner_pid, size, name, meta_binary, flag, creation_date, update_date, data_type, period FROM pretendo_smm.courses WHERE data_id=?`, dataID).Scan(&ownerPID, &size, &name, &metaBinary, &flag, &createdTime, &updatedTime, &dataType, &period)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil
		} else {
			log.Fatal(err)
		}
	}

	var stars uint32
	var attempts uint32
	var failures uint32
	var completions uint32

	_ = cassandraClusterSession.Query(`SELECT stars, attempts, failures, completions FROM pretendo_smm.ratings WHERE data_id=?`, dataID).Scan(&stars, &attempts, &failures, &completions)

	courseMetadata := &types.CourseMetadata{
		DataID:      dataID,
		OwnerPID:    ownerPID,
		Size:        size,
		CreatedTime: nex.NewDateTime(createdTime),
		UpdatedTime: nex.NewDateTime(updatedTime),
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

	return courseMetadata
}
