package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/utility"
)

func GetCustomRankingByDataID(err error, client *nex.Client, callID uint32, param *datastore_super_mario_maker_types.DataStoreGetCustomRankingByDataIDParam) uint32 {
	var pRankingResult []*datastore_super_mario_maker_types.DataStoreCustomRankingResult
	var pResults []uint32

	switch param.ApplicationID {
	case 0:
		if len(param.DataIDList) == 0 { // Starred courses
			pRankingResult, pResults = getCustomRankingByDataIdStarredCourses(client.PID())
		} else { // Played courses
			pRankingResult, pResults = getCustomRankingByDataIdCourseMetadata(param)
		}
	case 300000000: // Mii data
		pRankingResult, pResults = getCustomRankingByDataIdMiiData(param)
	default: // Normal metadata
		pRankingResult, pResults = getCustomRankingByDataIdCourseMetadata(param)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(pRankingResult)
	rmcResponseStream.WriteListUInt32LE(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetCustomRankingByDataID, rmcResponseBody)

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

	return 0
}

func getCustomRankingByDataIdStarredCourses(pid uint32) ([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, []uint32) {
	courseMetadatas := database.GetUserStarredCourses(pid)

	pRankingResult := make([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, 0)
	pResults := make([]uint32, 0)

	for _, courseMetadata := range courseMetadatas {
		pRankingResult = append(pRankingResult, utility.CourseMetadataToDataStoreCustomRankingResult(courseMetadata))
		pResults = append(pResults, 0x690001)
	}

	return pRankingResult, pResults
}

func getCustomRankingByDataIdMiiData(param *datastore_super_mario_maker_types.DataStoreGetCustomRankingByDataIDParam) ([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, []uint32) {
	pRankingResult := make([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, 0)
	pResults := make([]uint32, 0)

	for _, pid := range param.DataIDList {
		pid := uint32(pid)
		miiInfo := database.GetUserMiiInfoByPID(pid) // This isn't actually a PID when using the official servers! I set it as one to make this easier for me

		if miiInfo != nil {
			pRankingResult = append(pRankingResult, utility.UserMiiDataToDataStoreCustomRankingResult(pid, miiInfo))
			pResults = append(pResults, 0x690001)
		}
	}

	return pRankingResult, pResults
}

func getCustomRankingByDataIdCourseMetadata(param *datastore_super_mario_maker_types.DataStoreGetCustomRankingByDataIDParam) ([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, []uint32) {
	courseMetadatas := database.GetCourseMetadataByDataIDs(param.DataIDList)

	pRankingResult := make([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, 0)
	pResults := make([]uint32, 0)

	for _, courseMetadata := range courseMetadatas {
		pRankingResult = append(pRankingResult, utility.CourseMetadataToDataStoreCustomRankingResult(courseMetadata))
		pResults = append(pResults, 0x690001)
	}

	return pRankingResult, pResults
}
