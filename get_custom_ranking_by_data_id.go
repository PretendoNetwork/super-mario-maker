package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getCustomRankingByDataId(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreGetCustomRankingByDataIdParam) {
	var pRankingResult []*nexproto.DataStoreCustomRankingResult
	var pResults []uint32

	switch param.ApplicationId {
	case 0:
		if len(param.DataIdList) == 0 { // Starred courses
			pRankingResult, pResults = getCustomRankingByDataIdStarredCourses(client.PID())
		} else { // Played courses
			pRankingResult, pResults = getCustomRankingByDataIdCourseMetadata(param)
		}
	case 300000000: // Mii data
		pRankingResult, pResults = getCustomRankingByDataIdMiiData(param)
	default: // Normal metadata
		pRankingResult, pResults = getCustomRankingByDataIdCourseMetadata(param)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pRankingResult)
	rmcResponseStream.WriteListUInt32LE(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodGetCustomRankingByDataId, rmcResponseBody)

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

func getCustomRankingByDataIdStarredCourses(pid uint32) ([]*nexproto.DataStoreCustomRankingResult, []uint32) {
	courseMetadatas := getUserStarredCourses(pid)

	pRankingResult := make([]*nexproto.DataStoreCustomRankingResult, 0)
	pResults := make([]uint32, 0)

	for _, courseMetadata := range courseMetadatas {
		pRankingResult = append(pRankingResult, courseMetadataToDataStoreCustomRankingResult(courseMetadata))
		pResults = append(pResults, 0x690001)
	}

	return pRankingResult, pResults
}

func getCustomRankingByDataIdMiiData(param *nexproto.DataStoreGetCustomRankingByDataIdParam) ([]*nexproto.DataStoreCustomRankingResult, []uint32) {
	pRankingResult := make([]*nexproto.DataStoreCustomRankingResult, 0)
	pResults := make([]uint32, 0)

	for _, pid := range param.DataIdList {
		pid := uint32(pid)
		miiInfo := getUserMiiInfoByPID(pid) // This isn't actually a PID when using the official servers! I set it as one to make this easier for me

		if miiInfo != nil {
			pRankingResult = append(pRankingResult, userMiiDataToDataStoreCustomRankingResult(pid, miiInfo))
			pResults = append(pResults, 0x690001)
		}
	}

	return pRankingResult, pResults
}

func getCustomRankingByDataIdCourseMetadata(param *nexproto.DataStoreGetCustomRankingByDataIdParam) ([]*nexproto.DataStoreCustomRankingResult, []uint32) {
	courseMetadatas := getCourseMetadataByDataIDs(param.DataIdList)

	pRankingResult := make([]*nexproto.DataStoreCustomRankingResult, 0)
	pResults := make([]uint32, 0)

	for _, courseMetadata := range courseMetadatas {
		pRankingResult = append(pRankingResult, courseMetadataToDataStoreCustomRankingResult(courseMetadata))
		pResults = append(pResults, 0x690001)
	}

	return pRankingResult, pResults
}
