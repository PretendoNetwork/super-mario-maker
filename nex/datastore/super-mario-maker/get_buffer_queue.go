package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func GetBufferQueue(err error, packet nex.PacketInterface, callID uint32, param datastore_super_mario_maker_types.BufferQueueParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	pBufferQueue, nexError := datastore_smm_db.GetBufferQueuesByDataIDAndSlot(param.DataID, param.Slot)
	if nexError != nil {
		globals.Logger.Errorf("Error code %d for object %d", nexError.ResultCode, param.DataID)
		return nil, nexError
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pBufferQueue.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodGetBufferQueue
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
