package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getObjectInfos(err error, client *nex.Client, callID uint32, dataIDs []uint64) {
	pInfos := make([]*nexproto.DataStoreFileServerObjectInfo, 0)

	courseMetadatas := getCourseMetadataByDataIDs(dataIDs)

	for _, courseMetadata := range courseMetadatas {
		info := nexproto.NewDataStoreFileServerObjectInfo()
		info.DataID = courseMetadata.DataID
		info.GetInfo = nexproto.NewDataStoreReqGetInfo()
		info.GetInfo.URL = fmt.Sprintf("http://pds-AMAJ-d1.b-cdn.net/course/%d.bin", courseMetadata.DataID)
		info.GetInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
		info.GetInfo.Size = courseMetadata.Size
		info.GetInfo.RootCA = []byte{}
		info.GetInfo.DataID = courseMetadata.DataID

		pInfos = append(pInfos, info)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pInfos)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetObjectInfos, rmcResponseBody)

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
