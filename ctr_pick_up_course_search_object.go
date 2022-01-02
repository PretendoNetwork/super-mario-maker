package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

// This is the same as DataStoreSMM::RecommendedCourseSearchObject
// Not sure why they used a different method here?
func ctrPickUpCourseSearchObject(err error, client *nex.Client, callID uint32, dataStoreSearchParam *nexproto.DataStoreSearchParam, extraData []string) {
	// TODO complete this

	pRankingResults := make([]*nexproto.DataStoreCustomRankingResult, 0)

	courseMetadatas := getCourseMetadatasByLimit(100) // In PCAPs param.minimalRatingFrequency is 100 but is 0 here?

	for _, courseMetadata := range courseMetadatas {
		pRankingResults = append(pRankingResults, courseMetadataToDataStoreCustomRankingResult(courseMetadata))
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodCTRPickUpCourseSearchObject, rmcResponseBody)

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
