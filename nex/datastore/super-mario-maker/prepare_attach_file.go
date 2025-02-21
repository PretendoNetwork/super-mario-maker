package nex_datastore_super_mario_maker

import (
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func PrepareAttachFile(err error, packet nex.PacketInterface, callID uint32, param datastore_super_mario_maker_types.DataStoreAttachFileParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// * This method seems to be only used at "attach"
	// * a file to an existing object. In practice,
	// * SMM will use this to upload a courses preview
	// * image after uploading the course.
	// * param.ReferDataID is the courses object DataID

	dataID, nexError := datastore_smm_db.InitializeObjectByAttachFileParam(packet.Sender().PID(), param)
	if nexError != nil {
		globals.Logger.Errorf("Error code %d on object init", nexError.ResultCode)
		return nil, nexError
	}

	// TODO - Should this be moved to InitializeObjectByAttachFileParam?
	// * This never seems to have any values during normal gameplay,
	// * but just in case
	for i := range param.PostParam.RatingInitParams {
		nexError = datastore_db.InitializeObjectRatingWithSlot(uint64(dataID), param.PostParam.RatingInitParams[i])
		if nexError != nil {
			globals.Logger.Errorf("Error code %d on rating init", nexError.ResultCode)
			return nil, nexError
		}
	}

	// TODO - Check param.ContentType? Always seems to be "image/jpeg" but just in case?
	bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("%d.jpg", dataID)

	// TODO - Should this also take in the param.ContentType? To add it to the policy?
	URL, formData, _ := globals.Presigner.PostObject(bucket, key, time.Minute*15)

	pReqPostInfo := datastore_types.NewDataStoreReqPostInfo()

	pReqPostInfo.DataID = dataID
	pReqPostInfo.URL = types.NewString(URL.String())
	pReqPostInfo.FormFields = make(types.List[datastore_types.DataStoreKeyValue], 0, len(formData))

	for key, value := range formData {
		field := datastore_types.NewDataStoreKeyValue()
		field.Key = types.NewString(key)
		field.Value = types.NewString(value)

		pReqPostInfo.FormFields = append(pReqPostInfo.FormFields, field)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pReqPostInfo.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodPrepareAttachFile
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
