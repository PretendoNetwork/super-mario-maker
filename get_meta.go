package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMeta(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
	switch param.DataID {
	case 0: // Mii Data
		getMetaMiiData(client, callID, param)
	case 900000: // Course ID
		getMetaCourseMetadata(client, callID, param)
	default:
		fmt.Printf("[Warning] DataStoreProtocol::GetMeta Unsupported dataId: %v\n", param.DataID)
	}
}

func getMetaMiiData(client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
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

func getMetaCourseMetadata(client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaParam) {
	const examplePayload = "0062000000a0bb0d00000000000200000014de060001000032000000000500000000000000000005000000000000000086fcc87e1f0000001605a2861f00000032fb0000000000000000000000000086fcc87e1f00000000003e3f9c0000000000000000000000"
	examplePayloadBytes, _ := hex.DecodeString(examplePayload)

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMeta, examplePayloadBytes)

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
