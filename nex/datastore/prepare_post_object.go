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

func PreparePostObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStorePreparePostParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	dataID, errCode := datastore_db.InitializeObjectByPreparePostParam(client.PID(), param)
	if errCode != 0 {
		globals.Logger.Errorf("Error code %d on object init", errCode)
		return errCode
	}

	// TODO - Should this be moved to InitializeObjectByPreparePostParam?
	for _, ratingInitParamWithSlot := range param.RatingInitParams {
		errCode = datastore_db.InitializeObjectRatingWithSlot(dataID, ratingInitParamWithSlot)
		if errCode != 0 {
			globals.Logger.Errorf("Error code %d on rating init", errCode)
			return errCode
		}
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
