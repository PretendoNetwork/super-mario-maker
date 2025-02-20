package nex_datastore_super_mario_maker

import (
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func CompleteAttachFile(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreCompletePostParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// TODO - What is param.IsSuccess? Is this correct?
	if !param.IsSuccess {
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, "Invalid argument")
	}

	bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("%d.jpg", param.DataID)

	objectSizeS3, err := globals.S3ObjectSize(bucket, key)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
	}

	objectSizeDB, nexError := datastore_db.GetObjectSizeByDataID(param.DataID)
	if nexError != nil {
		return nil, nexError
	}

	if objectSizeS3 != uint64(objectSizeDB) {
		// TODO - Is this a good error?
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "")
	}

	nexError = datastore_db.UpdateObjectUploadCompletedByDataID(param.DataID, true)
	if nexError != nil {
		return nil, nexError
	}

	pURL, err := globals.Presigner.GetObject(bucket, key, time.Minute*15)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "Operation not allowed")
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	types.NewString(pURL.String()).WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodCompleteAttachFile
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
