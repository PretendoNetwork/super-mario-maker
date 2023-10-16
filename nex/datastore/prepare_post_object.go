package nex_datastore

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func PreparePostObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStorePreparePostParam) uint32 {
	rand.Seed(time.Now().UnixNano())
	nodeID := rand.Intn(len(globals.DataStoreIDGenerators))

	dataStoreIDGenerator := globals.DataStoreIDGenerators[nodeID]

	dataID := dataStoreIDGenerator.Next()
	database.SetDataStoreIDGeneratorLastID(nodeID, dataStoreIDGenerator.Value)
	database.InitializeCourseData(dataID, client.PID(), param.Size, param.Name, param.Flag, param.ExtraData, param.DataType, param.Period)

	if param.DataType != 1 { // 1 is Mii data, assume other values are course meta data
		database.UpdateCourseMetaBinary(dataID, param.MetaBinary)
	}

	bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("%d.bin", dataID)

	URL, formData, _ := globals.Presigner.PostObject(bucket, key, time.Minute*15)

	pReqPostInfo := datastore_types.NewDataStoreReqPostInfo()

	pReqPostInfo.DataID = dataID
	pReqPostInfo.URL = URL.String()
	pReqPostInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
	pReqPostInfo.FormFields = make([]*datastore_types.DataStoreKeyValue, 0, len(formData))
	pReqPostInfo.RootCACert = []byte{}

	for key, value := range formData {
		field := datastore_types.NewDataStoreKeyValue()
		field.Key = key
		field.Value = value

		pReqPostInfo.FormFields = append(pReqPostInfo.FormFields, field)
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(pReqPostInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPreparePostObject, rmcResponseBody)

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
