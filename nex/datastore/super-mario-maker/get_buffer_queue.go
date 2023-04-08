package nex_datastore_super_mario_maker

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetBufferQueue(err error, client *nex.Client, callID uint32, param *datastore_super_mario_maker.BufferQueueParam) {
	// TODO: complete this

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	var pBufferQueue [][]byte

	switch param.Slot {
	case 0: // unknown
		pBufferQueue = make([][]byte, 0)
	case 3: // death data
		pBufferQueue = database.GetBufferQueueDeathData(param.DataID)
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::GetBufferQueue Unsupported slot: %v\n", param.Slot)
	}

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

	globals.NEXServer.Send(responsePacket)
}
