package datastore_db

import (
	"time"

	"github.com/PretendoNetwork/nex-go"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/lib/pq"
)

func InitializeObjectByPreparePostParam(ownerPID uint32, param *datastore_types.DataStorePreparePostParam) (uint64, uint32) {
	now := time.Now()

	var dataID uint64

	err := database.Postgres.QueryRow(`INSERT INTO datastore.objects (
		owner,
		size,
		name,
		data_type,
		meta_binary,
		permission,
		permission_recipients,
		delete_permission,
		delete_permission_recipients,
		flag,
		period,
		refer_data_id,
		tags,
		persistence_slot_id,
		extra_data,
		creation_date,
		update_date
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17
	) RETURNING data_id`,
		ownerPID,
		param.Size,
		param.Name,
		param.DataType,
		param.MetaBinary,
		param.Permission.Permission,
		pq.Array(param.Permission.RecipientIDs),
		param.DelPermission.Permission,
		pq.Array(param.DelPermission.RecipientIDs),
		param.Flag,
		param.Period,
		param.ReferDataID,
		pq.Array(param.Tags),
		param.PersistenceInitParam.PersistenceSlotID, // TODO - Check param.PersistenceInitParam.DeleteLastObject?
		pq.Array(param.ExtraData),
		now,
		now,
	).Scan(&dataID)

	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return 0, nex.Errors.DataStore.Unknown
	}

	return dataID, 0
}
