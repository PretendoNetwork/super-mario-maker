package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func prepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *nexproto.DataStorePrepareGetParam) {
	// TODO: CDN

	pReqGetInfo := nexproto.NewDataStoreReqGetInfo()

	if dataStorePrepareGetParam.DataID == 900000 {
		pReqGetInfo.URL = "http://pds-AMAJ-d1.b-cdn.net/special/900000.bin"
		pReqGetInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
		pReqGetInfo.Size = 450068
		pReqGetInfo.RootCA = []byte{}
		pReqGetInfo.DataID = 900000
	} else {
		pReqGetInfo.URL = "http://pds-AMAJ-d1.b-cdn.net/course/1.bin"
		pReqGetInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
		pReqGetInfo.Size = 42516
		pReqGetInfo.RootCA = []byte{}
		pReqGetInfo.DataID = 1
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
