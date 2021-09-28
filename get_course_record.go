package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getCourseRecord(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreGetCourseRecordParam) {
	// TODO complete this
	// Hard coded to always say the course creator holds the world record

	/*
		courseMetadata := getCourseMetadataByDataID(param.DataID)

		result := nexproto.NewDataStoreGetCourseRecordResult()
		result.DataID = param.DataID
		result.Slot = 0
		result.FirstPID = courseMetadata.OwnerPID
		result.BestPID = courseMetadata.OwnerPID
		result.BestScore = 0
		result.CreatedTime = nex.NewDateTime(0)
		result.UpdatedTime = nex.NewDateTime(0)

		rmcResponseStream := nex.NewStreamOut(nexServer)

		rmcResponseStream.WriteStructure(result)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
		rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetCourseRecord, rmcResponseBody)

		rmcResponseBytes := rmcResponse.Bytes()
	*/

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetError(0x00690004)

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
