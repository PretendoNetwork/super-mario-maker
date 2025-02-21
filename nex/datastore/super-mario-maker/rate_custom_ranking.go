package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func RateCustomRanking(err error, packet nex.PacketInterface, callID uint32, params types.List[datastore_super_mario_maker_types.DataStoreRateCustomRankingParam]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// TODO - Check the period. The real server does check this, just unsure what it means or what the check is
	for i := range params {
		datastore_smm_db.InsertOrUpdateCustomRanking(params[i].DataID, params[i].ApplicationID, params[i].Score)
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodRateCustomRanking
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
