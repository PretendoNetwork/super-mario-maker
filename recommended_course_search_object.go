package main

import (
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func recommendedCourseSearchObject(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreSearchParam, extraData []string) {
	// TODO complete this

	pRankingResults := make([]*nexproto.DataStoreCustomRankingResult, 0)

	courseMetadatas := getCourseMetadatasByLimit(100) // In PCAPs param.minimalRatingFrequency is 100 but is 0 here?

	for i := 0; i < len(courseMetadatas); i++ {
		now := uint64(time.Now().Unix())
		courseMetadata := courseMetadatas[i]

		rankingResult := nexproto.NewDataStoreCustomRankingResult()

		rankingResult.Order = 0 // idk?
		rankingResult.Score = courseMetadata.Stars
		rankingResult.MetaInfo = nexproto.NewDataStoreMetaInfo()
		rankingResult.MetaInfo.DataID = courseMetadata.DataID
		rankingResult.MetaInfo.OwnerID = courseMetadata.OwnerPID
		rankingResult.MetaInfo.Size = courseMetadata.Size
		rankingResult.MetaInfo.Name = courseMetadata.Name
		rankingResult.MetaInfo.DataType = courseMetadata.DataType
		rankingResult.MetaInfo.MetaBinary = courseMetadata.MetaBinary
		rankingResult.MetaInfo.Permission = nexproto.NewDataStorePermission()
		rankingResult.MetaInfo.Permission.Permission = 0 // idk?
		rankingResult.MetaInfo.Permission.RecipientIds = []uint32{}
		rankingResult.MetaInfo.DelPermission = nexproto.NewDataStorePermission()
		rankingResult.MetaInfo.DelPermission.Permission = 3 // idk?
		rankingResult.MetaInfo.DelPermission.RecipientIds = []uint32{}
		rankingResult.MetaInfo.CreatedTime = nex.NewDateTime(now)
		rankingResult.MetaInfo.UpdatedTime = nex.NewDateTime(now)
		rankingResult.MetaInfo.Period = courseMetadata.Period
		rankingResult.MetaInfo.Status = 0      // idk?
		rankingResult.MetaInfo.ReferredCnt = 0 // idk?
		rankingResult.MetaInfo.ReferDataID = 0 // idk?
		rankingResult.MetaInfo.Flag = courseMetadata.Flag
		rankingResult.MetaInfo.ReferredTime = nex.NewDateTime(now)
		rankingResult.MetaInfo.ExpireTime = nex.NewDateTime(now)
		rankingResult.MetaInfo.Tags = []string{""} // idk?
		rankingResult.MetaInfo.Ratings = []*nexproto.DataStoreRatingInfoWithSlot{
			nexproto.NewDataStoreRatingInfoWithSlot(),
		}
		rankingResult.MetaInfo.Ratings[0].Slot = 0
		rankingResult.MetaInfo.Ratings[0].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[0].Rating.TotalValue = 3
		rankingResult.MetaInfo.Ratings[0].Rating.Count = 3
		rankingResult.MetaInfo.Ratings[0].Rating.InitialValue = 0

		pRankingResults = append(pRankingResults, rankingResult)
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
