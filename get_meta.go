package main

import (
	"encoding/hex"
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMeta(err error, client *nex.Client, callID uint32, dataStoreGetMetaParam *nexproto.DataStoreGetMetaParam) {
	switch dataStoreGetMetaParam.DataID {
	case 0: // Mii Data
		getMetaMiiData(client, callID, dataStoreGetMetaParam)
	case 900000: // Course ID
		getCourseMetadata(client, callID, dataStoreGetMetaParam)
	default:
		fmt.Printf("[Warning] DataStoreProtocol::GetMeta Unsupported dataId: %v\n", dataStoreGetMetaParam.DataID)
	}
}

func getMetaMiiData(client *nex.Client, callID uint32, dataStoreGetMetaParam *nexproto.DataStoreGetMetaParam) {
	const examplePayload = "00fc000000f86584030000000004395068000000000a005265644475636b73730001008c0042504643000000010000000000000000000000000001000003007330d94dae0fe5c6c34093aaf6a740f40769fe720000d2125200650064004400750063006b00730000000000642b000016010268441826344614811217688d008a25814948505200650064004400750063006b007300730000000000bf270000000000000000000000000000000000000001000500000000000000000005000000030000000073305f881f00000073305f881f0000005a000000000000000000000001000073305f881f00000000003e3f9c00000001000000030034390000000000"
	examplePayloadBytes, _ := hex.DecodeString(examplePayload)

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMeta, examplePayloadBytes)

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

func getCourseMetadata(client *nex.Client, callID uint32, dataStoreGetMetaParam *nexproto.DataStoreGetMetaParam) {
	const examplePayload = "0062000000a0bb0d00000000000200000014de060001000032000000000500000000000000000005000000000000000086fcc87e1f0000001605a2861f00000032fb0000000000000000000000000086fcc87e1f00000000003e3f9c0000000000000000000000"
	examplePayloadBytes, _ := hex.DecodeString(examplePayload)

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMeta, examplePayloadBytes)

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
