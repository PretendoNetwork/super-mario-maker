package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func CheckRateCustomRankingCounter(err error, packet nex.PacketInterface, callID uint32, applicationID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// * No idea what this is. Only seen application ID 0
	// * used, and it's always true? Unsure what this checks
	isBelowThreshold := types.NewBool(false)

	switch applicationID {
	case 0: // * Unknown
		isBelowThreshold = true
	default:
		globals.Logger.Warningf("Unsupported applicationID: %d", applicationID)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	isBelowThreshold.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodCheckRateCustomRankingCounter
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
