package main

import (
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func ctrPickUpCourseSearchObject(err error, client *nex.Client, callID uint32, dataStoreSearchParam *nexproto.DataStoreSearchParam, extraData []string) {
	// TODO complete this

	now := uint64(time.Now().Unix())

	pRankingResults := make([]*nexproto.DataStoreCustomRankingResult, 0)
	rankingResult := nexproto.NewDataStoreCustomRankingResult()

	rankingResult.Order = 0
	rankingResult.Score = 100
	rankingResult.MetaInfo = nexproto.NewDataStoreMetaInfo()
	rankingResult.MetaInfo.DataID = 1
	rankingResult.MetaInfo.OwnerID = 1730592963
	rankingResult.MetaInfo.Size = 42516
	rankingResult.MetaInfo.Name = "test"
	rankingResult.MetaInfo.DataType = 6 // idk?
	rankingResult.MetaInfo.MetaBinary = []byte{
		0x00, 0x00, 0x00, 0x01, // always 1
		0x00, 0x00, 0x00, 0x01, // course theme
		0x00, 0x00, 0x04, 0xcc, // length of compressed course data
		0x00, 0x00, 0x00, 0xc0, // length of compressed course sub data
		0x00, 0x00, 0x26, 0x8c, // length of compressed thumbnail0
		0x00, 0x00, 0x79, 0xfc, // length of compressed thumbnail1
		0x00, 0x00, 0x00, 0x01, // always 1
		0xcb, 0x38, 0x81, 0xaa, // crc32 of compressed course data
		0xdd, 0x9a, 0xda, 0xc3, // crc32 of compressed course sub data
		0xee, 0xf1, 0x18, 0x9f, // crc32 of compressed thumbnail0
		0x7a, 0x57, 0xf8, 0x71, // crc32 of compressed thumbnail1
	}
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
	rankingResult.MetaInfo.Flag = 0 // idk?
	rankingResult.MetaInfo.ReferredTime = nex.NewDateTime(now)
	rankingResult.MetaInfo.ExpireTime = nex.NewDateTime(now)
	rankingResult.MetaInfo.Tags = []string{""} // idk?
	rankingResult.MetaInfo.Ratings = []*nexproto.DataStoreRatingInfoWithSlot{}

	pRankingResults = append(pRankingResults, rankingResult)
	//0000000100000000000004cc000000c00000268c000079fc00000001cb3881aadd9adac3eef1189f7a57f871

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
