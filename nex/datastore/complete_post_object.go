package nex_datastore

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func CompletePostObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreCompletePostParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	// TODO - What is param.IsSuccess? Is this correct?
	if param.IsSuccess {
		bucket := os.Getenv("PN_SMM_CONFIG_S3_BUCKET")
		key := fmt.Sprintf("%d.bin", param.DataID)

		objectSizeS3, err := globals.S3ObjectSize(bucket, key)
		if err != nil {
			globals.Logger.Error(err.Error())
			return nex.Errors.DataStore.NotFound
		}

		objectSizeDB, errCode := datastore_db.GetObjectSizeDataID(param.DataID)
		if errCode != 0 {
			return errCode
		}

		if objectSizeS3 != uint64(objectSizeDB) {
			// TODO - Is this a good error?
			return nex.Errors.DataStore.Unknown
		}

		errCode = datastore_db.UpdateObjectUploadCompletedByDataID(param.DataID, true)
		if errCode != 0 {
			return errCode
		}
	}

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodCompletePostObject, nil)

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
