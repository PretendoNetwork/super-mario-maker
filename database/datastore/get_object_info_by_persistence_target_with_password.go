package datastore_db

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/lib/pq"
)

func GetObjectInfoByPersistenceTargetWithPassword(persistenceTarget *datastore_types.DataStorePersistenceTarget, password uint64) (*datastore_types.DataStoreMetaInfo, uint32) {
	metaInfo := datastore_types.NewDataStoreMetaInfo()
	metaInfo.Permission = datastore_types.NewDataStorePermission()
	metaInfo.DelPermission = datastore_types.NewDataStorePermission()
	metaInfo.CreatedTime = nex.NewDateTime(0)
	metaInfo.UpdatedTime = nex.NewDateTime(0)
	metaInfo.ReferredTime = nex.NewDateTime(0)
	metaInfo.ExpireTime = nex.NewDateTime(0x9C3f3E0000) // * 9999-12-31T00:00:00.000Z. This is what the real server sends
	metaInfo.Ratings = make([]*datastore_types.DataStoreRatingInfoWithSlot, 0)

	var accessPassword uint64
	var underReview bool
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
		update_date,
		access_password,
		under_review
	FROM datastore.objects WHERE owner=$1 AND persistence_slot_id=$2 AND upload_completed=TRUE AND deleted=FALSE`, persistenceTarget.OwnerID, persistenceTarget.PersistenceSlotID).Scan(
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
		&accessPassword,
		&underReview,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nex.Errors.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nil, nex.Errors.DataStore.Unknown
	}

	if accessPassword != 0 && accessPassword != password {
		return nil, nex.Errors.DataStore.InvalidPassword
	}

	if underReview {
		return nil, nex.Errors.DataStore.UnderReviewing
	}

	ratings, errCode := GetObjectRatingsWithSlotByDataIDWithPassword(metaInfo.DataID, password)
	if errCode != 0 {
		globals.Logger.Errorf("Failed to get ratings for object %d with password %d", metaInfo.DataID, password)
		return nil, errCode
	}

	metaInfo.Ratings = ratings

	metaInfo.CreatedTime.FromTimestamp(createdDate)
	metaInfo.UpdatedTime.FromTimestamp(updatedDate)
	metaInfo.ReferredTime.FromTimestamp(createdDate) // * This is what the real server does

	return metaInfo, 0
}
