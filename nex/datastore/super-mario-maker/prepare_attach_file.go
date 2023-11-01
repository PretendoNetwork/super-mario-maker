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
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func PrepareAttachFile(err error, packet nex.PacketInterface, callID uint32, param *datastore_super_mario_maker_types.DataStoreAttachFileParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	client := packet.Sender()

	// * This method seems to be only used at "attach"
	// * a file to an existing object. In practice,
	// * SMM will use this to upload a courses preview
	// * image after uploading the course.
	// * param.ReferDataID is the courses object DataID

	dataID, errCode := datastore_smm_db.InitializeObjectByAttachFileParam(client.PID(), param)
	if errCode != 0 {
		globals.Logger.Errorf("Error code %d on object init", errCode)
		return errCode
	}

	// TODO - Should this be moved to InitializeObjectByAttachFileParam?
	// * This never seems to have any values during normal gameplay,
	// * but just in case
	for _, ratingInitParamWithSlot := range param.PostParam.RatingInitParams {
		errCode = datastore_db.InitializeObjectRatingWithSlot(dataID, ratingInitParamWithSlot)
		if errCode != 0 {
			globals.Logger.Errorf("Error code %d on rating init", errCode)
			return errCode
		}
	}

	// TODO - Check param.ContentType? Always seems to be "image/jpeg" but just in case?
	bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("%d.jpg", dataID)

	// TODO - Should this also take in the param.ContentType? To add it to the policy?
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
