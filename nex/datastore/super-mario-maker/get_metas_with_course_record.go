package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetMetasWithCourseRecord(err error, packet nex.PacketInterface, callID uint32, params []*datastore_super_mario_maker_types.DataStoreGetCourseRecordParam, metaParam *datastore_types.DataStoreGetMetaParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	client := packet.Sender()

	// * The functionality here is largely a guess
	// * based on how GetCourseRecord works and
	// * based on testing with a custom client.
	// * During normal gameplay this method never
	// * seems to send any parameters, though if a
	// * custom client is used it IS functional

	pMetaInfo := make([]*datastore_types.DataStoreMetaInfo, 0, len(params))
	pCourseResults := make([]*datastore_super_mario_maker_types.DataStoreGetCourseRecordResult, 0, len(params))
	pResults := make([]*nex.Result, 0, len(params))

	for _, param := range params {
		// * metaParam has a password, but it's always set to 0.
		// * It also wouldn't make much sense for the same password
		// * to be used for all objects being requested here. So
		// * just assume metaParam is ONLY used for the resultOption
		// * field?
		objectInfo, errCode := datastore_db.GetObjectInfoByDataID(param.DataID)
		if errCode != 0 {
			// TODO - Maybe this should be broken out into a util function in globals?
			objectInfo = datastore_types.NewDataStoreMetaInfo()
			objectInfo.DataID = 0
			objectInfo.OwnerID = 0
			objectInfo.Size = 0
			objectInfo.Name = ""
			objectInfo.DataType = 0
			objectInfo.MetaBinary = []byte{}
			objectInfo.Permission = datastore_types.NewDataStorePermission()
			objectInfo.Permission.Permission = 0
			objectInfo.Permission.RecipientIDs = []uint32{}
			objectInfo.DelPermission = datastore_types.NewDataStorePermission()
			objectInfo.DelPermission.Permission = 0
			objectInfo.DelPermission.RecipientIDs = []uint32{}
			objectInfo.CreatedTime = nex.NewDateTime(0)
			objectInfo.UpdatedTime = nex.NewDateTime(0)
			objectInfo.Period = 0
			objectInfo.Status = 0
			objectInfo.ReferredCnt = 0
			objectInfo.ReferDataID = 0
			objectInfo.Flag = 0
			objectInfo.ReferredTime = nex.NewDateTime(0)
			objectInfo.ExpireTime = nex.NewDateTime(0)
			objectInfo.Tags = []string{}
			objectInfo.Ratings = []*datastore_types.DataStoreRatingInfoWithSlot{}
		} else {
			errCode = globals.VerifyObjectPermission(objectInfo.OwnerID, client.PID(), objectInfo.Permission)
			if errCode != 0 {
				objectInfo = datastore_types.NewDataStoreMetaInfo()
				objectInfo.DataID = 0
				objectInfo.OwnerID = 0
				objectInfo.Size = 0
				objectInfo.Name = ""
				objectInfo.DataType = 0
				objectInfo.MetaBinary = []byte{}
				objectInfo.Permission = datastore_types.NewDataStorePermission()
				objectInfo.Permission.Permission = 0
				objectInfo.Permission.RecipientIDs = []uint32{}
				objectInfo.DelPermission = datastore_types.NewDataStorePermission()
				objectInfo.DelPermission.Permission = 0
				objectInfo.DelPermission.RecipientIDs = []uint32{}
				objectInfo.CreatedTime = nex.NewDateTime(0)
				objectInfo.UpdatedTime = nex.NewDateTime(0)
				objectInfo.Period = 0
				objectInfo.Status = 0
				objectInfo.ReferredCnt = 0
				objectInfo.ReferDataID = 0
				objectInfo.Flag = 0
				objectInfo.ReferredTime = nex.NewDateTime(0)
				objectInfo.ExpireTime = nex.NewDateTime(0)
				objectInfo.Tags = []string{}
				objectInfo.Ratings = []*datastore_types.DataStoreRatingInfoWithSlot{}
			}

			// * This is kind of backwards.
			// * The database pulls this data
			// * by default, so it can be done
			// * in a single query. So instead
			// * of checking if a flag *IS*
			// * set, and conditionally *ADDING*
			// * the fields, we check if a flag
			// * is *NOT* set and conditionally
			// * *REMOVE* the field
			if metaParam.ResultOption&0x1 == 0 {
				objectInfo.Tags = make([]string, 0)
			}

			if metaParam.ResultOption&0x2 == 0 {
				objectInfo.Ratings = make([]*datastore_types.DataStoreRatingInfoWithSlot, 0)
			}

			if metaParam.ResultOption&0x4 == 0 {
				objectInfo.MetaBinary = make([]byte, 0)
			}
		}

		// * Ignore errors, real server sends empty struct if can't be found
		courseRecord, errCode := datastore_smm_db.GetCourseRecordByDataIDAndSlot(param.DataID, param.Slot)
		if errCode != 0 || objectInfo.DataID == 0 { // * DataID == 0 means could not be found or accessed
			courseRecord = datastore_super_mario_maker_types.NewDataStoreGetCourseRecordResult()
			courseRecord.DataID = 0
			courseRecord.Slot = 0
			courseRecord.FirstPID = 0
			courseRecord.BestPID = 0
			courseRecord.BestScore = 0
			courseRecord.CreatedTime = nex.NewDateTime(0)
			courseRecord.UpdatedTime = nex.NewDateTime(0)
		}

		pMetaInfo = append(pMetaInfo, objectInfo)
		pCourseResults = append(pCourseResults, courseRecord)
		pResults = append(pResults, nex.NewResultSuccess(nex.Errors.Core.Unknown)) // * Real server ALWAYS returns a success
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(pMetaInfo)
	rmcResponseStream.WriteListStructure(pCourseResults)
	rmcResponseStream.WriteListResult(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetMetasWithCourseRecord, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
