package nex_datastore_super_mario_maker

import (
	"fmt"
	"os"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetObjectInfos(err error, client *nex.Client, callID uint32, dataIDs []uint64) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	pInfos := make([]*datastore_super_mario_maker_types.DataStoreFileServerObjectInfo, 0)

	for _, dataID := range dataIDs {
		objectInfo, errCode := datastore_db.GetObjectInfoByDataID(dataID)
		if errCode != 0 {
			return errCode
		}

		bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
		key := fmt.Sprintf("%d.bin", objectInfo.DataID)

		URL, err := globals.Presigner.GetObject(bucket, key, time.Minute*15)
		if err != nil {
			globals.Logger.Error(err.Error())
			return nex.Errors.DataStore.OperationNotAllowed
		}

		info := datastore_super_mario_maker_types.NewDataStoreFileServerObjectInfo()
		info.DataID = objectInfo.DataID
		info.GetInfo = datastore_types.NewDataStoreReqGetInfo()
		info.GetInfo.URL = URL.String()
		info.GetInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
		info.GetInfo.Size = objectInfo.Size
		info.GetInfo.RootCACert = []byte{}
		info.GetInfo.DataID = objectInfo.DataID

		pInfos = append(pInfos, info)
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(pInfos)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetObjectInfos, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
