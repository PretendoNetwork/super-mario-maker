package nex_datastore

import (
	"fmt"
	"os"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/utility"
)

func PrepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *datastore_types.DataStorePrepareGetParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.Core.Unknown
	}

	pReqGetInfo := datastore_types.NewDataStoreReqGetInfo()

	bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("%d.bin", dataStorePrepareGetParam.DataID)

	objectSize, err := utility.S3ObjectSize(bucket, key)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.NotFound
	}

	URL, err := globals.Presigner.GetObject(bucket, key, time.Minute*15)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.OperationNotAllowed
	}

	pReqGetInfo.URL = URL.String()
	pReqGetInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
	pReqGetInfo.Size = uint32(objectSize)
	pReqGetInfo.RootCACert = []byte{}
	pReqGetInfo.DataID = dataStorePrepareGetParam.DataID

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

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

	globals.NEXServer.Send(responsePacket)

	return 0
}
