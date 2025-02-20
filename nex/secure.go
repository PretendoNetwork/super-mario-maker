package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func StartSecureServer() {
	serverBuildString = "branch:origin/project/nfs build:3_10_26_2006_0"

	globals.SecureServer = nex.NewPRUDPServer()

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = false

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 8, 3))
	globals.SecureServer.AccessKey = "9f2b4678"

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== SMM1 - Secure ===")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("==================")
	})

	// * Register the common handlers first so that they can be overridden if needed
	registerCommonSecureProtocols()
	registerNEXProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_SMM_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
