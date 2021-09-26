package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func completeAttachFile(err error, client *nex.Client, callID uint32, dataStoreCompletePostParam *nexproto.DataStoreCompletePostParam) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO complete this

	/*
		fmt.Println(dataStoreCompletePostParam.IsSuccess)
		fmt.Println(dataStoreCompletePostParam.DataID)

		if dataStoreCompletePostParam.IsSuccess {
			setCoursePlayable(dataStoreCompletePostParam.DataID)
		}
	*/

	rmcResponseStream.WriteString(fmt.Sprintf("http://pds-amaj-d1.b-cdn.net/image/%d.jpg", dataStoreCompletePostParam.DataID))

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodCompleteAttachFile, rmcResponseBody)

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
