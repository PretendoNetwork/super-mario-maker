package nex_datastore

import (
	"fmt"
	"os"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func PrepareGetObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStorePrepareGetParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.Core.Unknown
	}

	bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("%d.bin", param.DataID)

	objectInfo, errCode := datastore_db.GetObjectInfoByDataID(param.DataID)
	if errCode != 0 {
		return errCode
	}

	errCode = globals.VerifyObjectPermission(objectInfo.OwnerID, client.PID(), objectInfo.Permission)
	if errCode != 0 {
		return errCode
	}

	URL, err := globals.Presigner.GetObject(bucket, key, time.Minute*15)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.OperationNotAllowed
	}

	pReqGetInfo := datastore_types.NewDataStoreReqGetInfo()

	pReqGetInfo.URL = URL.String()
	pReqGetInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
	pReqGetInfo.Size = objectInfo.Size
	pReqGetInfo.RootCACert = []byte{}
	pReqGetInfo.DataID = param.DataID

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(pReqGetInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPrepareGetObject, rmcResponseBody)

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
