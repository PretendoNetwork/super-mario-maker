package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func GetMetasWithCourseRecord(err error, packet nex.PacketInterface, callID uint32, params types.List[datastore_super_mario_maker_types.DataStoreGetCourseRecordParam], metaParam datastore_types.DataStoreGetMetaParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// * The functionality here is largely a guess
	// * based on how GetCourseRecord works and
	// * based on testing with a custom client.
	// * During normal gameplay this method never
	// * seems to send any parameters, though if a
	// * custom client is used it IS functional

	pMetaInfo := make(types.List[datastore_types.DataStoreMetaInfo], 0, len(params))
	pCourseResults := make(types.List[datastore_super_mario_maker_types.DataStoreGetCourseRecordResult], 0, len(params))
	pResults := make(types.List[types.QResult], 0, len(params))

	for _, param := range params {
		// * metaParam has a password, but it's always set to 0.
		// * It also wouldn't make much sense for the same password
		// * to be used for all objects being requested here. So
		// * just assume metaParam is ONLY used for the resultOption
		// * field?
		objectInfo, nexError := datastore_db.GetObjectInfoByDataID(param.DataID)
		if nexError != nil {
			objectInfo = datastore_types.NewDataStoreMetaInfo()
		} else {
			nexError = globals.DatastoreCommon.VerifyObjectPermission(objectInfo.OwnerID, packet.Sender().PID(), objectInfo.Permission)
			if nexError != nil {
				objectInfo = datastore_types.NewDataStoreMetaInfo()
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
				objectInfo.Tags = types.NewList[types.String]()
			}

			if metaParam.ResultOption&0x2 == 0 {
				objectInfo.Ratings = types.NewList[datastore_types.DataStoreRatingInfoWithSlot]()
			}

			if metaParam.ResultOption&0x4 == 0 {
				objectInfo.MetaBinary = types.NewQBuffer(nil)
			}
		}

		// * Ignore errors, real server sends empty struct if can't be found
		courseRecord, nexError := datastore_smm_db.GetCourseRecordByDataIDAndSlot(param.DataID, param.Slot)
		if nexError != nil || objectInfo.DataID == 0 { // * DataID == 0 means could not be found or accessed
			courseRecord = datastore_super_mario_maker_types.NewDataStoreGetCourseRecordResult()
		}

		pMetaInfo = append(pMetaInfo, objectInfo)
		pCourseResults = append(pCourseResults, courseRecord)
		pResults = append(pResults, types.NewQResultSuccess(nex.ResultCodes.Core.Unknown)) // * Real server ALWAYS returns a success
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pMetaInfo.WriteTo(rmcResponseStream)
	pCourseResults.WriteTo(rmcResponseStream)
	pResults.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodGetMetasWithCourseRecord
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
