package nex

import (
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	nex_secure_connection_common "github.com/PretendoNetwork/super-mario-maker-secure/nex/secure-connection/common"
)

func registerCommonProtocols() {
	secureConnectionProtocol := secureconnection.NewCommonSecureConnectionProtocol(globals.NEXServer)

	secureConnectionProtocol.AddConnection(nex_secure_connection_common.AddConnection)             // * Stubbed
	secureConnectionProtocol.UpdateConnection(nex_secure_connection_common.UpdateConnection)       // * Stubbed
	secureConnectionProtocol.DoesConnectionExist(nex_secure_connection_common.DoesConnectionExist) // * Stubbed
}
