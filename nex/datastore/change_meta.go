package nex_datastore

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func ChangeMeta(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreChangeMetaParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	metaInfo, errCode := datastore_db.GetObjectInfoByDataID(param.DataID)
	if errCode != 0 {
		return errCode
	}

	// TODO - Is this the right permission?
	errCode = globals.VerifyObjectPermission(metaInfo.OwnerID, client.PID(), metaInfo.DelPermission)
	if errCode != 0 {
		return errCode
	}

	// TODO - Refactor this into a single query. Use a query builder like goqu?
	if param.ModifiesFlag&0x08 != 0 {
		errCode = datastore_db.UpdateObjectPeriodByDataIDWithPassword(param.DataID, param.Period, param.UpdatePassword)
		if errCode != 0 {
			return errCode
		}
	}

	if param.ModifiesFlag&0x10 != 0 {
		errCode = datastore_db.UpdateObjectMetaBinaryByDataIDWithPassword(param.DataID, param.MetaBinary, param.UpdatePassword)
		if errCode != 0 {
			return errCode
		}
	}

	if param.ModifiesFlag&0x80 != 0 {
		errCode = datastore_db.UpdateObjectDataTypeByDataIDWithPassword(param.DataID, param.DataType, param.UpdatePassword)
		if errCode != 0 {
			return errCode
		}
	}

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodChangeMeta, nil)

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
