package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func prepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *nexproto.DataStorePrepareGetParam) {
	pReqGetInfo := nexproto.NewDataStoreReqGetInfo()

	if dataStorePrepareGetParam.DataID == 900000 {
		pReqGetInfo.URL = "http://pds-AMAJ-d1.b-cdn.net/special/900000.bin"
		pReqGetInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
		pReqGetInfo.Size = 450068
		pReqGetInfo.RootCA = []byte{}
		pReqGetInfo.DataID = 900000
	} else {
		courseMetadata := getCourseMetadataByDataID(dataStorePrepareGetParam.DataID)

		pReqGetInfo.URL = fmt.Sprintf("http://pds-AMAJ-d1.b-cdn.net/course/%d.bin", dataStorePrepareGetParam.DataID)
		pReqGetInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
		pReqGetInfo.Size = courseMetadata.Size
		pReqGetInfo.RootCA = []byte{}
		pReqGetInfo.DataID = dataStorePrepareGetParam.DataID
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteStructure(pReqGetInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodPrepareGetObject, rmcResponseBody)

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
