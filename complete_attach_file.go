package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func completeAttachFile(err error, client *nex.Client, callID uint32, dataStoreCompletePostParam *nexproto.DataStoreCompletePostParam) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO: complete this

	rmcResponseStream.WriteString(fmt.Sprintf("http://%s.b-cdn.net/image/%d.jpg", os.Getenv("DO_SPACES_NAME"), dataStoreCompletePostParam.DataID))

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
