package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func UploadCourseRecord(err error, packet nex.PacketInterface, callID uint32, param datastore_super_mario_maker_types.DataStoreUploadCourseRecordParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	client := packet.Sender()

	nexError := datastore_smm_db.InsertOrUpdateCourseRecord(param.DataID, param.Slot, client.PID(), param.Score)
	if nexError != nil {
		return nil, nexError
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodUploadCourseRecord
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
