package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetCourseRecord(err error, client *nex.Client, callID uint32, param *datastore_super_mario_maker.DataStoreGetCourseRecordParam) {
	worldRecord := database.GetCourseWorldRecord(param.DataID)

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)

	if worldRecord == nil {
		rmcResponse.SetError(0x690004)
	} else {
		result := datastore_super_mario_maker.NewDataStoreGetCourseRecordResult()
		result.DataID = param.DataID
		result.Slot = param.Slot
		result.FirstPID = worldRecord.FirstPID
		result.BestPID = worldRecord.BestPID
		result.BestScore = worldRecord.Score
		result.CreatedTime = worldRecord.CreatedTime
		result.UpdatedTime = worldRecord.UpdatedTime

		rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

		rmcResponseStream.WriteStructure(result)

		rmcResponseBody := rmcResponseStream.Bytes()
		rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetCourseRecord, rmcResponseBody)
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

	globals.NEXServer.Send(responsePacket)
}
