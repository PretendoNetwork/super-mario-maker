package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetBufferQueue(err error, client *nex.Client, callID uint32, param *datastore_super_mario_maker_types.BufferQueueParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	pBufferQueue, errCode := datastore_smm_db.GetBufferQueuesByDataIDAndSlot(param.DataID, param.Slot)
	if errCode != 0 {
		globals.Logger.Errorf("Error code %d for object %d", errCode, param.DataID)
		return errCode
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)
	rmcResponseStream.WriteListQBuffer(pBufferQueue)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetBufferQueue, rmcResponseBody)

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
