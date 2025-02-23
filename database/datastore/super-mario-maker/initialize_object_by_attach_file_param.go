package datastore_smm_db

import (
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_smm_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	"github.com/PretendoNetwork/super-mario-maker/database"
	"github.com/PretendoNetwork/super-mario-maker/globals"
	"github.com/lib/pq"
)

func InitializeObjectByAttachFileParam(ownerPID types.PID, param datastore_smm_types.DataStoreAttachFileParam) (types.UInt64, *nex.Error) {
	now := time.Now()

	var dataID types.UInt64
	tagArray := make([]string, 0, len(param.PostParam.Tags))
	for i := range param.PostParam.Tags {
		tagArray = append(tagArray, string(param.PostParam.Tags[i]))
	}

	extraDataArray := make([]string, 0, len(param.PostParam.ExtraData))
	for i := range param.PostParam.Tags {
		extraDataArray = append(extraDataArray, string(param.PostParam.ExtraData[i]))
	}

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
		param.PostParam.Size,
		param.PostParam.Name,
		param.PostParam.DataType,
		param.PostParam.MetaBinary,
		param.PostParam.Permission.Permission,
		pq.Array(param.PostParam.Permission.RecipientIDs),
		param.PostParam.DelPermission.Permission,
		pq.Array(param.PostParam.DelPermission.RecipientIDs),
		param.PostParam.Flag,
		param.PostParam.Period,
		param.ReferDataID,
		pq.Array(tagArray),
		param.PostParam.PersistenceInitParam.PersistenceSlotID, // TODO - Check param.PersistenceInitParam.DeleteLastObject?
		pq.Array(extraDataArray),
		now,
		now,
	).Scan(&dataID)

	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return types.NewUInt64(0), nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return dataID, nil
}
