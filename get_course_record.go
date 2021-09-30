package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getCourseRecord(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreGetCourseRecordParam) {
	worldRecord := getCourseWorldRecord(param.DataID)

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)

	if worldRecord == nil {
		rmcResponse.SetError(0x00690004)
	} else {
		result := nexproto.NewDataStoreGetCourseRecordResult()
		result.DataID = param.DataID
		result.Slot = param.Slot
		result.FirstPID = worldRecord.FirstPID
		result.BestPID = worldRecord.BestPID
		result.BestScore = worldRecord.Score
		result.CreatedTime = worldRecord.CreatedTime
		result.UpdatedTime = worldRecord.UpdatedTime

		rmcResponseStream := nex.NewStreamOut(nexServer)

		rmcResponseStream.WriteStructure(result)

		rmcResponseBody := rmcResponseStream.Bytes()
		rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetCourseRecord, rmcResponseBody)
	}

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
