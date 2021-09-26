package main

import (
	"encoding/base64"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getCustomRankingByDataId(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreGetCustomRankingByDataIdParam) {
	if param.ApplicationId == 300000000 {
		// Mii Data
		getCustomRankingByDataIdMiiData(client, callID, param)
	} else {
		// Course Metadata
		getCustomRankingByDataIdCourseMetadata(client, callID, param)
	}
}

func getCustomRankingByDataIdMiiData(client *nex.Client, callID uint32, param *nexproto.DataStoreGetCustomRankingByDataIdParam) {
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

func getCustomRankingByDataIdCourseMetadata(client *nex.Client, callID uint32, param *nexproto.DataStoreGetCustomRankingByDataIdParam) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO complete this

	rmcResponseStream.WriteUInt32LE(0x00000000) // pRankingResult List length 0
	rmcResponseStream.WriteUInt32LE(0x00000000) // pResults List length 0

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
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
