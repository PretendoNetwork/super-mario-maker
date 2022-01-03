package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMetasMultipleParam(err error, client *nex.Client, callID uint32, params []*nexproto.DataStoreGetMetaParam) {
	pMetaInfo := make([]*nexproto.DataStoreMetaInfo, 0)
	pResults := make([]uint32, 0)

	for _, param := range params {
		if param.DataID == 0 {
			pMetaInfo = append(pMetaInfo, getMetasMultipleParamMiiData(param))
		} else {
			fmt.Println("Unknown meta multiple data ID")
		}

		pResults = append(pResults, 0x690001)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pMetaInfo)
	rmcResponseStream.WriteListUInt32LE(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMetasMultipleParam, rmcResponseBody)

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

func getMetasMultipleParamMiiData(param *nexproto.DataStoreGetMetaParam) *nexproto.DataStoreMetaInfo {
	miiInfo := getUserMiiInfoByPID(param.PersistenceTarget.OwnerID)

	return userMiiDataToDataStoreMetaInfo(param.PersistenceTarget.OwnerID, miiInfo)
}
