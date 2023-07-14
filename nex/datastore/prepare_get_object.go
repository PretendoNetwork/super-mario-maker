package nex_datastore

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/utility"
)

func PrepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *datastore_types.DataStorePrepareGetParam) {
	pReqGetInfo := datastore_types.NewDataStoreReqGetInfo()

	if dataStorePrepareGetParam.DataID == 900000 {
		objectSize, _ := utility.S3ObjectSize(os.Getenv("S3_BUCKET_NAME"), "special/900000.bin")

		pReqGetInfo.URL = fmt.Sprintf("http://%s.b-cdn.net/special/900000.bin", "pds-AMAJ-d1")
		pReqGetInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
		pReqGetInfo.Size = uint32(objectSize)
		pReqGetInfo.RootCACert = []byte{}
		pReqGetInfo.DataID = 900000
	} else {
		courseMetadata := database.GetCourseMetadataByDataID(dataStorePrepareGetParam.DataID)

		pReqGetInfo.URL = fmt.Sprintf("http://%s.b-cdn.net/course/%d.bin", os.Getenv("S3_BUCKET_NAME"), dataStorePrepareGetParam.DataID)
		pReqGetInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
		pReqGetInfo.Size = courseMetadata.Size
		pReqGetInfo.RootCACert = []byte{}
		pReqGetInfo.DataID = dataStorePrepareGetParam.DataID
	}

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
}
