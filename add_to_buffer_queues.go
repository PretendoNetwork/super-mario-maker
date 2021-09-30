package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func addToBufferQueues(err error, client *nex.Client, callID uint32, params []*nexproto.BufferQueueParam, buffers [][]byte) {
	pResults := make([]uint32, 0)

	for i := 0; i < len(params); i++ {
		buffer := buffers[i]
		param := params[i]

		insertBufferQueueData(param.DataID, param.Slot, buffer)

		pResults = append(pResults, 0x690001)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListUInt32LE(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodAddToBufferQueues, rmcResponseBody)

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
