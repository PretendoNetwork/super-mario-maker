package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMetasWithCourseRecord(err error, client *nex.Client, callID uint32, dataStoreGetCourseRecordParams []*nexproto.DataStoreGetCourseRecordParam, dataStoreGetMetaParam *nexproto.DataStoreGetMetaParam) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO complete this

	rmcResponseStream.WriteUInt32LE(0x00000000) // pMetaInfo List length 0
	rmcResponseStream.WriteUInt32LE(0x00000000) // pCourseResults List length 0
	rmcResponseStream.WriteUInt32LE(0x00000000) // pResults List length 0

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetMetasWithCourseRecord, rmcResponseBody)

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
