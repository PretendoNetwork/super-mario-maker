package nex_datastore_super_mario_maker

import (
	"fmt"
	"os"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func PrepareAttachFile(err error, client *nex.Client, callID uint32, dataStoreAttachFileParam *datastore_super_mario_maker_types.DataStoreAttachFileParam) uint32 {
	dataID := dataStoreAttachFileParam.ReferDataID
	bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("%d.jpeg", dataID)

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

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodPrepareAttachFile, rmcResponseBody)

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
