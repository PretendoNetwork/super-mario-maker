package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func changeMeta(err error, client *nex.Client, callID uint32, dataStoreChangeMetaParam *nexproto.DataStoreChangeMetaParam) {
	// STUBBED
	// TODO: DO SOMETHING WITH THE DATA
	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodChangeMeta, nil)

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
