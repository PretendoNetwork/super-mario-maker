package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func recommendedCourseSearchObject(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreSearchParam, extraData []string) {
	// TODO: complete this

	pRankingResults := make([]*nexproto.DataStoreCustomRankingResult, 0)

	// TEMP FOR SHUTTER TO TEST THINGS
	if client.PID() == 1049991375 {
		courseMetadata := getCourseMetadataByDataID(145) // specific course shutter wants
		pRankingResults = append(pRankingResults, courseMetadataToDataStoreCustomRankingResult(courseMetadata))
	} else {
		courseMetadatas := getCourseMetadatasByLimit(100) // In PCAPs param.minimalRatingFrequency is 100 but is 0 here?

		for _, courseMetadata := range courseMetadatas {
			pRankingResults = append(pRankingResults, courseMetadataToDataStoreCustomRankingResult(courseMetadata))
		}
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodRecommendedCourseSearchObject, rmcResponseBody)

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
