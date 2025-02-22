package datastore_db

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
	"github.com/lib/pq"
)

func GetObjectInfoByDataID(dataID types.UInt64) (datastore_types.DataStoreMetaInfo, *nex.Error) {
	nexError := IsObjectAvailable(dataID)
	if nexError != nil {
		return datastore_types.NewDataStoreMetaInfo(), nexError
	}

	metaInfo := datastore_types.NewDataStoreMetaInfo()
	metaInfo.Permission = datastore_types.NewDataStorePermission()
	metaInfo.DelPermission = datastore_types.NewDataStorePermission()
	metaInfo.ExpireTime = types.NewDateTime(0x9C3F3E0000) // * 9999-12-31T00:00:00.000Z. This is what the real server sends
	metaInfo.Ratings = types.NewList[datastore_types.DataStoreRatingInfoWithSlot]()

	var createdDate time.Time
	var updatedDate time.Time

	err := database.Postgres.QueryRow(`SELECT
		data_id,
		owner,
		size,
		name,
		data_type,
		meta_binary,
		permission,
		permission_recipients,
		delete_permission,
		delete_permission_recipients,
		period,
		refer_data_id,
		flag,
		tags,
		creation_date,
		update_date
	FROM datastore.objects WHERE data_id=$1`, dataID).Scan(
		&metaInfo.DataID,
		&metaInfo.OwnerID,
		&metaInfo.Size,
		&metaInfo.Name,
		&metaInfo.DataType,
		&metaInfo.MetaBinary,
		&metaInfo.Permission.Permission,
		pq.Array(&metaInfo.Permission.RecipientIDs),
		&metaInfo.DelPermission.Permission,
		pq.Array(&metaInfo.DelPermission.RecipientIDs),
		&metaInfo.Period,
		&metaInfo.ReferDataID,
		&metaInfo.Flag,
		pq.Array(&metaInfo.Tags),
		&createdDate,
		&updatedDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return datastore_types.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return datastore_types.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	ratings, nexError := GetObjectRatingsWithSlotByDataID(metaInfo.DataID)
	if nexError != nil {
		return datastore_types.NewDataStoreMetaInfo(), nexError
	}

	metaInfo.Ratings = ratings

	metaInfo.CreatedTime.FromTimestamp(createdDate)
	metaInfo.UpdatedTime.FromTimestamp(updatedDate)
	metaInfo.ReferredTime.FromTimestamp(createdDate) // * This is what the real server does

	return metaInfo, nil
}
