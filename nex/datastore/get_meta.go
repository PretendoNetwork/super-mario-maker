package nex_datastore

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/utility"
)

func GetMeta(err error, client *nex.Client, callID uint32, param *datastore.DataStoreGetMetaParam) {
	switch param.DataID {
	case 0: // Mii Data
		getMetaMiiData(client, callID, param)
	case 900000: // Event course news
		getMetaEventCourseNewsData(client, callID, param)
	default:
		fmt.Printf("[Warning] DataStoreProtocol::GetMeta Unsupported dataId: %v\n", param.DataID)
	}
}

func getMetaMiiData(client *nex.Client, callID uint32, param *datastore.DataStoreGetMetaParam) {
	miiInfo := database.GetUserMiiInfoByPID(param.PersistenceTarget.OwnerID)

	pMetaInfo := utility.UserMiiDataToDataStoreMetaInfo(param.PersistenceTarget.OwnerID, miiInfo)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)
	rmcResponseStream.WriteStructure(pMetaInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodGetMeta, rmcResponseBody)

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
}

func getMetaEventCourseNewsData(client *nex.Client, callID uint32, param *datastore.DataStoreGetMetaParam) {
	objectSize, _ := utility.S3ObjectSize(os.Getenv("S3_BUCKET_NAME"), "special/900000.bin")

	pMetaInfo := datastore.NewDataStoreMetaInfo()
	pMetaInfo.DataID = 900000
	pMetaInfo.OwnerID = 2
	pMetaInfo.Size = uint32(objectSize)
	pMetaInfo.Name = ""
	pMetaInfo.DataType = 50 // Metdata?
	pMetaInfo.MetaBinary = []byte{}
	pMetaInfo.Permission = datastore.NewDataStorePermission()
	pMetaInfo.Permission.Permission = 0 // idk?
	pMetaInfo.Permission.RecipientIds = []uint32{}
	pMetaInfo.DelPermission = datastore.NewDataStorePermission()
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
	pMetaInfo.Ratings = []*datastore.DataStoreRatingInfoWithSlot{}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)
	rmcResponseStream.WriteStructure(pMetaInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodGetMeta, rmcResponseBody)

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
}
