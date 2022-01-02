package main

import (
	"encoding/base64"
	"time"

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
		rankingResult.MetaInfo.CreatedTime = courseMetadata.CreatedTime
		rankingResult.MetaInfo.UpdatedTime = courseMetadata.UpdatedTime
		rankingResult.MetaInfo.Period = courseMetadata.Period
		rankingResult.MetaInfo.Status = 0      // unknown
		rankingResult.MetaInfo.ReferredCnt = 0 // unknown
		rankingResult.MetaInfo.ReferDataID = 0 // unknown
		rankingResult.MetaInfo.Flag = courseMetadata.Flag
		rankingResult.MetaInfo.ReferredTime = courseMetadata.CreatedTime
		rankingResult.MetaInfo.ExpireTime = nex.NewDateTime(671075926016) // December 31st, year 9999
		rankingResult.MetaInfo.Tags = []string{""}                        // unknown
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
		rankingResult.MetaInfo.Ratings[0].Rating.InitialValue = 0

		// unknown
		rankingResult.MetaInfo.Ratings[1].Slot = 1
		rankingResult.MetaInfo.Ratings[1].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[1].Rating.TotalValue = 2
		rankingResult.MetaInfo.Ratings[1].Rating.Count = 2
		rankingResult.MetaInfo.Ratings[1].Rating.InitialValue = 0

		// completions
		rankingResult.MetaInfo.Ratings[2].Slot = 2
		rankingResult.MetaInfo.Ratings[2].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[2].Rating.TotalValue = int64(courseMetadata.Completions)
		rankingResult.MetaInfo.Ratings[2].Rating.Count = courseMetadata.Completions
		rankingResult.MetaInfo.Ratings[2].Rating.InitialValue = 0

		// failures
		rankingResult.MetaInfo.Ratings[3].Slot = 3
		rankingResult.MetaInfo.Ratings[3].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[3].Rating.TotalValue = int64(courseMetadata.Failures)
		rankingResult.MetaInfo.Ratings[3].Rating.Count = courseMetadata.Failures
		rankingResult.MetaInfo.Ratings[3].Rating.InitialValue = 0

		// unknown
		rankingResult.MetaInfo.Ratings[4].Slot = 4
		rankingResult.MetaInfo.Ratings[4].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[4].Rating.TotalValue = 5
		rankingResult.MetaInfo.Ratings[4].Rating.Count = 5
		rankingResult.MetaInfo.Ratings[4].Rating.InitialValue = 0

		// unknown
		rankingResult.MetaInfo.Ratings[5].Slot = 5
		rankingResult.MetaInfo.Ratings[5].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[5].Rating.TotalValue = 6
		rankingResult.MetaInfo.Ratings[5].Rating.Count = 6
		rankingResult.MetaInfo.Ratings[5].Rating.InitialValue = 0

		// Number of new Miiverse comments
		rankingResult.MetaInfo.Ratings[6].Slot = 6
		rankingResult.MetaInfo.Ratings[6].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[6].Rating.TotalValue = 0
		rankingResult.MetaInfo.Ratings[6].Rating.Count = 0
		rankingResult.MetaInfo.Ratings[6].Rating.InitialValue = 0

		pRankingResult = append(pRankingResult, rankingResult)
		pResults = append(pResults, 0x690001)
	}

	return pRankingResult, pResults
}

func getCustomRankingByDataIdMiiData(param *nexproto.DataStoreGetCustomRankingByDataIdParam) ([]*nexproto.DataStoreCustomRankingResult, []uint32) {
	pRankingResult := make([]*nexproto.DataStoreCustomRankingResult, 0)
	pResults := make([]uint32, 0)

	for i := 0; i < len(param.DataIdList); i++ {
		ownerID := uint32(param.DataIdList[i]) // This isn't actually a PID when using the official servers! I set it as one to make this easier for me
		miiInfo := getUserMiiInfoByPID(ownerID)

		if miiInfo != nil {
			encodedMiiData := miiInfo["data"].(string)
			decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

			metaBinaryStream := nex.NewStreamOut(nexServer)
			metaBinaryStream.Grow(140)
			metaBinaryStream.WriteBytesNext([]byte{
				0x42, 0x50, 0x46, 0x43, // BPFC magic
				0x00, 0x00, 0x00, 0x01, // Unknown
				0x00, 0x00, 0x00, 0x00, // Unknown
				0x00, 0x00, 0x00, 0x00, // Unknown
				0x00, 0x00, 0x00, 0x00, // Unknown
				0x00, 0x01, 0x00, 0x00, // Unknown
			})
			metaBinaryStream.WriteBytesNext(decodedMiiData) // Actual Mii data
			metaBinaryStream.WriteBytesNext([]byte{
				0x00, 0x00, 0x00, 0x00, // Unknown
				0x00, 0x00, 0x00, 0x00, // Unknown
				0x00, 0x00, 0x00, 0x00, // Unknown
				0x00, 0x00, 0x00, 0x00, // Unknown
				0x00, 0x00, 0x00, 0x01, // Unknown
			})

			now := uint64(time.Now().Unix())

			rankingResult := nexproto.NewDataStoreCustomRankingResult()

			rankingResult.Order = 0
			rankingResult.Score = 0
			rankingResult.MetaInfo = nexproto.NewDataStoreMetaInfo()
			rankingResult.MetaInfo.DataID = uint64(ownerID) // idk what this is, but it gets used elsewhere for request Mii data again. Setting it as a PID makes that easier for me
			rankingResult.MetaInfo.OwnerID = ownerID
			rankingResult.MetaInfo.Size = 0
			rankingResult.MetaInfo.Name = miiInfo["name"].(string)
			rankingResult.MetaInfo.DataType = 1 // Mii data type?
			rankingResult.MetaInfo.MetaBinary = metaBinaryStream.Bytes()
			rankingResult.MetaInfo.Permission = nexproto.NewDataStorePermission()
			rankingResult.MetaInfo.Permission.Permission = 0 // idk?
			rankingResult.MetaInfo.Permission.RecipientIds = []uint32{}
			rankingResult.MetaInfo.DelPermission = nexproto.NewDataStorePermission()
			rankingResult.MetaInfo.DelPermission.Permission = 3 // idk?
			rankingResult.MetaInfo.DelPermission.RecipientIds = []uint32{}
			rankingResult.MetaInfo.CreatedTime = nex.NewDateTime(now)
			rankingResult.MetaInfo.UpdatedTime = nex.NewDateTime(now)
			rankingResult.MetaInfo.Period = 90 // idk?
			rankingResult.MetaInfo.Status = 0
			rankingResult.MetaInfo.ReferredCnt = 0
			rankingResult.MetaInfo.ReferDataID = 0
			rankingResult.MetaInfo.Flag = 256 // idk?
			rankingResult.MetaInfo.ReferredTime = nex.NewDateTime(now)
			rankingResult.MetaInfo.ExpireTime = nex.NewDateTime(now)
			rankingResult.MetaInfo.Tags = []string{"49"} // idk?
			rankingResult.MetaInfo.Ratings = []*nexproto.DataStoreRatingInfoWithSlot{}

			pRankingResult = append(pRankingResult, rankingResult)
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
		rankingResult.MetaInfo.CreatedTime = courseMetadata.CreatedTime
		rankingResult.MetaInfo.UpdatedTime = courseMetadata.UpdatedTime
		rankingResult.MetaInfo.Period = courseMetadata.Period
		rankingResult.MetaInfo.Status = 0      // unknown
		rankingResult.MetaInfo.ReferredCnt = 0 // unknown
		rankingResult.MetaInfo.ReferDataID = 0 // unknown
		rankingResult.MetaInfo.Flag = courseMetadata.Flag
		rankingResult.MetaInfo.ReferredTime = courseMetadata.CreatedTime
		rankingResult.MetaInfo.ExpireTime = nex.NewDateTime(671075926016) // December 31st, year 9999
		rankingResult.MetaInfo.Tags = []string{""}                        // unknown
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
		rankingResult.MetaInfo.Ratings[0].Rating.InitialValue = 0

		// unknown
		rankingResult.MetaInfo.Ratings[1].Slot = 1
		rankingResult.MetaInfo.Ratings[1].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[1].Rating.TotalValue = 2
		rankingResult.MetaInfo.Ratings[1].Rating.Count = 2
		rankingResult.MetaInfo.Ratings[1].Rating.InitialValue = 0

		// completions
		rankingResult.MetaInfo.Ratings[2].Slot = 2
		rankingResult.MetaInfo.Ratings[2].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[2].Rating.TotalValue = int64(courseMetadata.Completions)
		rankingResult.MetaInfo.Ratings[2].Rating.Count = courseMetadata.Completions
		rankingResult.MetaInfo.Ratings[2].Rating.InitialValue = 0

		// failures
		rankingResult.MetaInfo.Ratings[3].Slot = 3
		rankingResult.MetaInfo.Ratings[3].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[3].Rating.TotalValue = int64(courseMetadata.Failures)
		rankingResult.MetaInfo.Ratings[3].Rating.Count = courseMetadata.Failures
		rankingResult.MetaInfo.Ratings[3].Rating.InitialValue = 0

		// unknown
		rankingResult.MetaInfo.Ratings[4].Slot = 4
		rankingResult.MetaInfo.Ratings[4].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[4].Rating.TotalValue = 5
		rankingResult.MetaInfo.Ratings[4].Rating.Count = 5
		rankingResult.MetaInfo.Ratings[4].Rating.InitialValue = 0

		// unknown
		rankingResult.MetaInfo.Ratings[5].Slot = 5
		rankingResult.MetaInfo.Ratings[5].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[5].Rating.TotalValue = 6
		rankingResult.MetaInfo.Ratings[5].Rating.Count = 6
		rankingResult.MetaInfo.Ratings[5].Rating.InitialValue = 0

		// Number of new Miiverse comments
		rankingResult.MetaInfo.Ratings[6].Slot = 6
		rankingResult.MetaInfo.Ratings[6].Rating = nexproto.NewDataStoreRatingInfo()
		rankingResult.MetaInfo.Ratings[6].Rating.TotalValue = 0
		rankingResult.MetaInfo.Ratings[6].Rating.Count = 0
		rankingResult.MetaInfo.Ratings[6].Rating.InitialValue = 0

		pRankingResult = append(pRankingResult, rankingResult)
		pResults = append(pResults, 0x690001)
	}

	return pRankingResult, pResults
}
