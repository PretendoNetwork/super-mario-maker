package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMeta(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
	switch param.DataID {
	case 0: // Mii Data
		getMetaMiiData(client, callID, param)
	case 900000: // Event course news
		getMetaEventCourseNewsData(client, callID, param)
	default:
		fmt.Printf("[Warning] DataStoreProtocol::GetMeta Unsupported dataId: %v\n", param.DataID)
	}
}

func getMetaMiiData(client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
	miiInfo := getUserMiiInfoByPID(param.PersistenceTarget.OwnerID)

	pMetaInfo := userMiiDataToDataStoreMetaInfo(param.PersistenceTarget.OwnerID, miiInfo)

	rmcResponseStream := nex.NewStreamOut(nexServer)
	rmcResponseStream.WriteStructure(pMetaInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMeta, rmcResponseBody)

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

func getMetaEventCourseNewsData(client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
	objectSize, _ := s3ObjectSize(os.Getenv("S3_BUCKET_NAME"), "special/900000.bin")

	pMetaInfo := nexproto.NewDataStoreMetaInfo()
	pMetaInfo.DataID = 900000
	pMetaInfo.OwnerID = 2
	pMetaInfo.Size = uint32(objectSize)
	pMetaInfo.Name = ""
	pMetaInfo.DataType = 50 // Metdata?
	pMetaInfo.MetaBinary = []byte{}
	pMetaInfo.Permission = nexproto.NewDataStorePermission()
	pMetaInfo.Permission.Permission = 0 // idk?
	pMetaInfo.Permission.RecipientIds = []uint32{}
	pMetaInfo.DelPermission = nexproto.NewDataStorePermission()
	pMetaInfo.DelPermission.Permission = 0 // idk?
	pMetaInfo.DelPermission.RecipientIds = []uint32{}
	pMetaInfo.CreatedTime = nex.NewDateTime(135271087238) // Reused from Nintendo
	pMetaInfo.UpdatedTime = nex.NewDateTime(135402751254) // Reused from Nintendo
	pMetaInfo.Period = 64306                              // idk?
	pMetaInfo.Status = 0
	pMetaInfo.ReferredCnt = 0
	pMetaInfo.ReferDataID = 0
	pMetaInfo.Flag = 0                                     // idk?
	pMetaInfo.ReferredTime = nex.NewDateTime(135271087238) // Reused from Nintendo
	pMetaInfo.ExpireTime = nex.NewDateTime(671075926016)   // Reused from Nintendo
	pMetaInfo.Tags = []string{}                            // idk?
	pMetaInfo.Ratings = []*nexproto.DataStoreRatingInfoWithSlot{}

	rmcResponseStream := nex.NewStreamOut(nexServer)
	rmcResponseStream.WriteStructure(pMetaInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMeta, rmcResponseBody)

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
