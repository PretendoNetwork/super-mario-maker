package main

import (
	"encoding/base64"
	"fmt"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMetasMultipleParam(err error, client *nex.Client, callID uint32, params []*nexproto.DataStoreGetMetaParam) {
	pMetaInfo := make([]*nexproto.DataStoreMetaInfo, 0)
	pResults := make([]uint32, 0)

	for i := 0; i < len(params); i++ {
		param := params[i]

		if param.DataID == 0 {
			pMetaInfo = append(pMetaInfo, getMetasMultipleParamMiiData(param))
		} else {
			fmt.Println("Unknown meta multiple data ID")
		}

		pResults = append(pResults, 0x690001)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pMetaInfo)
	rmcResponseStream.WriteListUInt32LE(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMetasMultipleParam, rmcResponseBody)

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

func getMetasMultipleParamMiiData(param *nexproto.DataStoreGetMetaParam) *nexproto.DataStoreMetaInfo {
	miiInfo := getUserMiiInfoByPID(param.PersistenceTarget.OwnerID)

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

	pMetaInfo := nexproto.NewDataStoreMetaInfo()

	pMetaInfo.DataID = uint64(param.PersistenceTarget.OwnerID) // idk what this is, but it gets used elsewhere for request Mii data again. Setting it as a PID makes that easier for me
	pMetaInfo.OwnerID = param.PersistenceTarget.OwnerID
	pMetaInfo.Size = 0
	pMetaInfo.Name = miiInfo["name"].(string)
	pMetaInfo.DataType = 1 // Mii data type?
	pMetaInfo.MetaBinary = metaBinaryStream.Bytes()
	pMetaInfo.Permission = nexproto.NewDataStorePermission()
	pMetaInfo.Permission.Permission = 0 // idk?
	pMetaInfo.Permission.RecipientIds = []uint32{}
	pMetaInfo.DelPermission = nexproto.NewDataStorePermission()
	pMetaInfo.DelPermission.Permission = 3 // idk?
	pMetaInfo.DelPermission.RecipientIds = []uint32{}
	pMetaInfo.CreatedTime = nex.NewDateTime(now)
	pMetaInfo.UpdatedTime = nex.NewDateTime(now)
	pMetaInfo.Period = 90 // idk?
	pMetaInfo.Status = 0
	pMetaInfo.ReferredCnt = 0
	pMetaInfo.ReferDataID = 0
	pMetaInfo.Flag = 256 // idk?
	pMetaInfo.ReferredTime = nex.NewDateTime(now)
	pMetaInfo.ExpireTime = nex.NewDateTime(now)
	pMetaInfo.Tags = []string{"49"} // idk?
	pMetaInfo.Ratings = []*nexproto.DataStoreRatingInfoWithSlot{}

	return pMetaInfo
}
