package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func AddToBufferQueues(err error, packet nex.PacketInterface, callID uint32, params types.List[datastore_super_mario_maker_types.BufferQueueParam], buffers types.List[types.QBuffer]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	client := packet.Sender()

	pResults := make(types.List[types.QResult], 0)

	// * The number of params and buffers CAN be allowed
	// * to differ, though this doesn't appear to happen
	// * in normal gameplay. In these cases, the real
	// * server will just loop over smallest list and
	// * ignore all other data in the larger list
	iterations := min(len(params), len(buffers))

	for i := 0; i < iterations; i++ {
		param := params[i]
		buffer := buffers[i]

		if param.Slot == 0 {
			objectInfo, nexError := datastore_db.GetObjectInfoByDataID(param.DataID)
			if nexError != nil {
				return nil, nexError
			}

			// * Objects with DataType 1 are "maker" objects. When adding
			// * a buffer to slot 0 of a maker object, a course is being
			// * added to that users "Starred Courses" list. The buffer
			// * is the courses DataID saved as an 8 byte buffer. To
			// * prevent people from adding courses to random users
			// * "Starred Courses" lists, we have to verify the requesting
			// * client owns the maker object
			if objectInfo.DataType == 1 && objectInfo.OwnerID != client.PID() {
				return nil, nex.NewError(nex.ResultCodes.DataStore.PermissionDenied, "Permission denied")
			}
		}

		nexError := datastore_smm_db.InsertOrUpdateBufferQueueData(param.DataID, param.Slot, buffer)
		if nexError != nil {
			return nil, nexError
		}

		pResults = append(pResults, types.NewQResultSuccess(nex.ResultCodes.Core.Unknown)) // * Seems to ALWAYS be a success?
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pResults.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodAddToBufferQueues
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
