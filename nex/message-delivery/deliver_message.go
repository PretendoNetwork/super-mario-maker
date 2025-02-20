package nex_message_delivery

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	message_delivery "github.com/PretendoNetwork/nex-protocols-go/v2/message-delivery"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func DeliverMessage(err error, packet nex.PacketInterface, callID uint32, oUserMessage types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// TODO - See what this does

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = message_delivery.ProtocolID
	rmcResponse.MethodID = message_delivery.MethodDeliverMessage
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
