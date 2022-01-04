package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func prepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *nexproto.DataStorePrepareGetParam) {
	pReqGetInfo := nexproto.NewDataStoreReqGetInfo()

	if dataStorePrepareGetParam.DataID == 900000 {
		pReqGetInfo.URL = fmt.Sprintf("http://%s.b-cdn.net/special/900000.bin", os.Getenv("DO_SPACES_NAME"))
		pReqGetInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
		//pReqGetInfo.Size = 450068
		pReqGetInfo.Size = 0x263C
		pReqGetInfo.RootCA = []byte{}
		pReqGetInfo.DataID = 900000
	} else {
		courseMetadata := getCourseMetadataByDataID(dataStorePrepareGetParam.DataID)

		pReqGetInfo.URL = fmt.Sprintf("http://%s.b-cdn.net/course/%d.bin", os.Getenv("DO_SPACES_NAME"), dataStorePrepareGetParam.DataID)
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
