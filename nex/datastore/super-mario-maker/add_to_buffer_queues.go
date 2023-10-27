package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func AddToBufferQueues(err error, client *nex.Client, callID uint32, params []*datastore_super_mario_maker_types.BufferQueueParam, buffers [][]byte) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	pResults := make([]*nex.Result, 0)

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
			objectInfo, errCode := datastore_db.GetObjectInfoByDataID(param.DataID)
			if errCode != 0 {
				return errCode
			}

			// * Objects with DataType 1 are "maker" objects. When adding
			// * a buffer to slot 0 of a maker object, a course is being
			// * added to that users "Starred Courses" list. The buffer
			// * is the courses DataID saved as an 8 byte buffer. To
			// * prevent people from adding courses to random users
			// * "Starred Courses" lists, we have to verify the requesting
			// * client owns the maker object
			if objectInfo.DataType == 1 && objectInfo.OwnerID != client.PID() {
				return nex.Errors.DataStore.PermissionDenied
			}
		}

		errCode := datastore_smm_db.InsertOrUpdateBufferQueueData(param.DataID, param.Slot, buffer)
		if errCode != 0 {
			return errCode
		}

		pResults = append(pResults, nex.NewResultSuccess(nex.Errors.Core.Unknown)) // * Seems to ALWAYS be a success?
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListResult(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodAddToBufferQueues, rmcResponseBody)

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
