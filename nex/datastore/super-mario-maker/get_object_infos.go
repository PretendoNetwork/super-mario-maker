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
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func GetObjectInfos(err error, packet nex.PacketInterface, callID uint32, dataIDs types.List[types.UInt64]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	pInfos := types.NewList[datastore_super_mario_maker_types.DataStoreFileServerObjectInfo]()

	for _, dataID := range dataIDs {
		objectInfo, nexError := datastore_db.GetObjectInfoByDataID(dataID)
		if nexError != nil {
			return nil, nexError
		}

		bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
		key := fmt.Sprintf("%d.bin", objectInfo.DataID)

		URL, err := globals.Presigner.GetObject(bucket, key, time.Minute*15)
		if err != nil {
			globals.Logger.Error(err.Error())
			return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "Operation not allowed")
		}

		info := datastore_super_mario_maker_types.NewDataStoreFileServerObjectInfo()
		info.DataID = objectInfo.DataID
		info.GetInfo = datastore_types.NewDataStoreReqGetInfo()
		info.GetInfo.URL = types.NewString(URL.String())
		info.GetInfo.Size = objectInfo.Size
		info.GetInfo.DataID = objectInfo.DataID

		pInfos = append(pInfos, info)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pInfos.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodGetObjectInfos
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
