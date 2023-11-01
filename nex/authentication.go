package nex

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

var serverBuildString string

func StartAuthenticationServer() {
	globals.AuthenticationServer = nex.NewServer()
	globals.AuthenticationServer.SetPRUDPVersion(1)
	globals.AuthenticationServer.SetPRUDPProtocolMinorVersion(2)
	globals.AuthenticationServer.SetDefaultNEXVersion(nex.NewPatchedNEXVersion(3, 8, 3, "AMAJ"))
	globals.AuthenticationServer.SetKerberosPassword(globals.KerberosPassword)
	globals.AuthenticationServer.SetAccessKey("9f2b4678")

	globals.AuthenticationServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==SMM1 - Auth==")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID())
		fmt.Printf("Method ID: %d\n", request.MethodID())
		fmt.Println("===============")
	})

	serverBuildString = "build:3_8_29_3022_0"

	registerCommonAuthenticationServerProtocols()

	globals.AuthenticationServer.Listen(fmt.Sprintf(":%s", os.Getenv("PN_SMM_AUTHENTICATION_SERVER_PORT")))
}
