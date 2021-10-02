package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func checkRateCustomRankingCounter(err error, client *nex.Client, callID uint32, applicationID uint32) {

	var isBelowThreshold uint8

	switch applicationID {
	case 0: // Unknown?
		isBelowThreshold = 1
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::CheckRateCustomRankingCounter Unsupported applicationID: %v\n", applicationID)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt8(isBelowThreshold)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodCheckRateCustomRankingCounter, rmcResponseBody)

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
