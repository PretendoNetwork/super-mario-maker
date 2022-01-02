package main

import (
	"encoding/hex"
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMeta(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
	switch param.DataID {
	case 0: // Mii Data
		getMetaMiiData(client, callID, param)
	case 900000: // Course ID
		getMetaCourseMetadata(client, callID, param)
	default:
		fmt.Printf("[Warning] DataStoreProtocol::GetMeta Unsupported dataId: %v\n", param.DataID)
	}
}

func getMetaMiiData(client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
	miiInfo := getUserMiiInfoByPID(param.PersistenceTarget.OwnerID)

	pMetaInfo := userMiiDataToDataStoreMetaInfo(param.PersistenceTarget.OwnerID, miiInfo)

	rmcResponseStream := nex.NewStreamOut(nexServer)
	rmcResponseStream.WriteStructure(pMetaInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMeta, rmcResponseBody)

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

func getMetaCourseMetadata(client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
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
