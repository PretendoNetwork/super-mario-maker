package nex

import (
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func registerCommonSecureProtocols() {
	secureconnection.NewCommonSecureConnectionProtocol(globals.SecureServer)
}
