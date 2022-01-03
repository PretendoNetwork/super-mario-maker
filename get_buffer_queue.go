package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getBufferQueue(err error, client *nex.Client, callID uint32, param *nexproto.BufferQueueParam) {
	// TODO: complete this

	rmcResponseStream := nex.NewStreamOut(nexServer)

	var pBufferQueue [][]byte

	switch param.Slot {
	case 0: // unknown
		pBufferQueue = make([][]byte, 0)
	case 3: // death data
		pBufferQueue = getBufferQueueDeathData(param.DataID)
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::GetBufferQueue Unsupported slot: %v\n", param.Slot)
	}

	rmcResponseStream.WriteListQBuffer(pBufferQueue)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetBufferQueue, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
