package nex_datastore_super_mario_maker

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

var MAX_COURSE_UPLOADS uint32 = 0x64 // Nintendo has this set to 100, but we can change it if we want

func GetApplicationConfig(err error, client *nex.Client, callID uint32, applicationID uint32) uint32 {

	config := make([]uint32, 0)

	switch applicationID {
	case 0: // Unknown?
		config = getApplicationConfig_Unknown0()
	case 1: // Official player PIDs
		config = getApplicationConfig_OfficialPlayers()
	case 2: // Unknown?
		config = getApplicationConfig_Unknown2()
	case 10: // Unknown?
		config = getApplicationConfig_Unknown10()
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::GetApplicationConfig Unsupported applicationID: %v\n", applicationID)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListUInt32LE(config)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetApplicationConfig, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)

	return 0
}

func getApplicationConfig_Unknown0() []uint32 {
	// I have no idea what this is
	// Just replaying data sent from the real server
	return []uint32{
		0x01000000, 0x32000000, 0x96000000, 0x2c010000, 0xf4010000,
		0x20030000, 0x14050000, 0xd0070000, 0xb80b0000, 0x88130000,
		MAX_COURSE_UPLOADS, 0x14000000, 0x1e000000, 0x28000000, 0x32000000,
		0x3c000000, 0x46000000, 0x50000000, 0x5a000000, 0x64000000,
		0x23000000, 0x4b000000, 0x23000000, 0x4b000000, 0x32000000,
		0x00000000, 0x03000000, 0x03000000, 0x64000000, 0x06000000,
		0x01000000, 0x60000000, 0x05000000, 0x60000000, 0x00000000,
		0xe4070000, 0x01000000, 0x01000000, 0x0c000000, 0x00000000,
	}
}

func getApplicationConfig_OfficialPlayers() []uint32 {
	return []uint32{
		0x02000000, // 2 (not a real user PID, this translates to the internal Quazal Rendez-Vous user used by NEX)
		0x70cc8269, // 1770179696 (official_player0 on NN, need to make PN versions)
		0x50cc8269, // 1770179664 (official_player1 on NN, need to make PN versions)
		0x38cc8269, // 1770179640 (official_player2 on NN, need to make PN versions)
		0xdbd08269, // 1770180827 (official_player3 on NN, need to make PN versions)
		0xa9d08269, // 1770180777 (official_player4 on NN, need to make PN versions)
		0x89d08269, // 1770180745 (official_player5 on NN, need to make PN versions)
		0x59c48269, // 1770177625 (official_player6 on NN, need to make PN versions)
		0x36c48269, // 1770177590 (official_player7 on NN, need to make PN versions)
	}
}

func getApplicationConfig_Unknown2() []uint32 {
	// I have no idea what this is
	// Just replaying data sent from the real server
	return []uint32{0xdf070000, 0x0c000000, 0x16000000, 0x05000000, 0x00000000}
}

func getApplicationConfig_Unknown10() []uint32 {
	// I have no idea what this is
	// Just replaying data sent from the real server
	// Only seen on the 3DS
	return []uint32{35, 75, 96, 40, 5, 6}
}
