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

		rankingResult.Order = 0 // unknown
		rankingResult.Score = courseMetadata.Stars
		rankingResult.MetaInfo = nexproto.NewDataStoreMetaInfo()
		rankingResult.MetaInfo.DataID = courseMetadata.DataID
		rankingResult.MetaInfo.OwnerID = courseMetadata.OwnerPID
		rankingResult.MetaInfo.Size = courseMetadata.Size
		rankingResult.MetaInfo.Name = courseMetadata.Name
		rankingResult.MetaInfo.DataType = courseMetadata.DataType
		rankingResult.MetaInfo.MetaBinary = courseMetadata.MetaBinary
		rankingResult.MetaInfo.Permission = nexproto.NewDataStorePermission()
		rankingResult.MetaInfo.Permission.Permission = 0 // unknown
		rankingResult.MetaInfo.Permission.RecipientIds = []uint32{}
		rankingResult.MetaInfo.DelPermission = nexproto.NewDataStorePermission()
		rankingResult.MetaInfo.DelPermission.Permission = 3 // unknown
		rankingResult.MetaInfo.DelPermission.RecipientIds = []uint32{}
		rankingResult.MetaInfo.CreatedTime = nex.NewDateTime(now)
		rankingResult.MetaInfo.UpdatedTime = nex.NewDateTime(now)
		rankingResult.MetaInfo.Period = courseMetadata.Period
		rankingResult.MetaInfo.Status = 0      // unknown
		rankingResult.MetaInfo.ReferredCnt = 0 // unknown
		rankingResult.MetaInfo.ReferDataID = 0 // unknown
		rankingResult.MetaInfo.Flag = courseMetadata.Flag
		rankingResult.MetaInfo.ReferredTime = nex.NewDateTime(now)
		rankingResult.MetaInfo.ExpireTime = nex.NewDateTime(now)
		rankingResult.MetaInfo.Tags = []string{""} // unknown
		rankingResult.MetaInfo.Ratings = []*nexproto.DataStoreRatingInfoWithSlot{
			nexproto.NewDataStoreRatingInfoWithSlot(), // attempts
			nexproto.NewDataStoreRatingInfoWithSlot(), // unknown
			nexproto.NewDataStoreRatingInfoWithSlot(), // completions
			nexproto.NewDataStoreRatingInfoWithSlot(), // failures
			nexproto.NewDataStoreRatingInfoWithSlot(), // unknown
			nexproto.NewDataStoreRatingInfoWithSlot(), // unknown
			nexproto.NewDataStoreRatingInfoWithSlot(), // unknown
		}

		// attempts
		rankingResult.MetaInfo.Ratings[0].Slot = 0
		rankingResult.MetaInfo.Ratings[0].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[0].Rating.TotalValue = int64(courseMetadata.Attempts)
		rankingResult.MetaInfo.Ratings[0].Rating.Count = courseMetadata.Attempts
		rankingResult.MetaInfo.Ratings[0].Rating.InitialValue = int64(courseMetadata.Attempts)

		// unknown
		rankingResult.MetaInfo.Ratings[1].Slot = 1
		rankingResult.MetaInfo.Ratings[1].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[1].Rating.TotalValue = 0
		rankingResult.MetaInfo.Ratings[1].Rating.Count = 0
		rankingResult.MetaInfo.Ratings[1].Rating.InitialValue = 0

		// completions
		rankingResult.MetaInfo.Ratings[2].Slot = 2
		rankingResult.MetaInfo.Ratings[2].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[2].Rating.TotalValue = int64(courseMetadata.Completions)
		rankingResult.MetaInfo.Ratings[2].Rating.Count = courseMetadata.Completions
		rankingResult.MetaInfo.Ratings[2].Rating.InitialValue = int64(courseMetadata.Completions)

		// failures
		rankingResult.MetaInfo.Ratings[3].Slot = 3
		rankingResult.MetaInfo.Ratings[3].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[3].Rating.TotalValue = int64(courseMetadata.Failures)
		rankingResult.MetaInfo.Ratings[3].Rating.Count = courseMetadata.Failures
		rankingResult.MetaInfo.Ratings[3].Rating.InitialValue = int64(courseMetadata.Failures)

		// unknown
		rankingResult.MetaInfo.Ratings[4].Slot = 4
		rankingResult.MetaInfo.Ratings[4].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[4].Rating.TotalValue = 5
		rankingResult.MetaInfo.Ratings[4].Rating.Count = 5
		rankingResult.MetaInfo.Ratings[4].Rating.InitialValue = 5

		// unknown
		rankingResult.MetaInfo.Ratings[5].Slot = 5
		rankingResult.MetaInfo.Ratings[5].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[5].Rating.TotalValue = 6
		rankingResult.MetaInfo.Ratings[5].Rating.Count = 6
		rankingResult.MetaInfo.Ratings[5].Rating.InitialValue = 6

		// unknown
		rankingResult.MetaInfo.Ratings[6].Slot = 6
		rankingResult.MetaInfo.Ratings[6].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[6].Rating.TotalValue = 7
		rankingResult.MetaInfo.Ratings[6].Rating.Count = 7
		rankingResult.MetaInfo.Ratings[6].Rating.InitialValue = 7

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
