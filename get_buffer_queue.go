package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getBufferQueue(err error, client *nex.Client, callID uint32, param *nexproto.BufferQueueParam) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO: complete this

	var pBufferQueue [][]byte

	switch param.Slot {
	case 0: // unknown
		pBufferQueue = make([][]byte, 0)
	case 3: // death data
		pBufferQueue = getBufferQueueDeathData(param.DataID)
		incrementCourseAttemptCount(param.DataID) // We also know this is when a user attempts a course
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
