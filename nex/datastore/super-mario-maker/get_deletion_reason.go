package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func GetDeletionReason(err error, packet nex.PacketInterface, callID uint32, dataIDLst types.List[types.UInt64]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// TODO - Complete this
	// * It's not actually known what the
	// * real "deletion reason" values are.
	// * This is stubbed until we figure
	// * that out
	pDeletionReasons := types.NewList[types.UInt32]()

	for range dataIDLst {
		// * Every course I've checked has had this
		// * set to 0, even if the course is not
		// * deleted
		pDeletionReasons = append(pDeletionReasons, types.NewUInt32(0))
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pDeletionReasons.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodGetDeletionReason
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
