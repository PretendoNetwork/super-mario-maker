package nex

import (
	message_delivery "github.com/PretendoNetwork/nex-protocols-go/message-delivery"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	nex_message_delivery "github.com/PretendoNetwork/super-mario-maker-secure/nex/message-delivery"
)

func registerNEXProtocols() {
	messageDeliveryProtocol := message_delivery.NewProtocol(globals.SecureServer)

	messageDeliveryProtocol.DeliverMessage(nex_message_delivery.DeliverMessage)
}
