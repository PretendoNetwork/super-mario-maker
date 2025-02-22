package nex_datastore_super_mario_maker

import (
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

// * Nintendo sets this to 10 by default
// * and users earn more upload slots up
// * to 100.
// * This is a stupid, unfun, mechanic so
// * everyone gets 100 by default. Can be
// * more, but 100 is fine tbh
var MAX_COURSE_UPLOADS uint32 = 100

func GetApplicationConfig(err error, packet nex.PacketInterface, callID uint32, applicationID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	var config []uint32

	switch applicationID {
	case 0: // * Player config?
		config = getApplicationConfig_PlayerConfig()
	case 1: // * PIDs of the "Official" makers in the "MAKERS" section
		config = getApplicationConfig_OfficialMakers()
	case 2: // * Unknown
		config = getApplicationConfig_Unknown2()
	case 10: // * Unknown
		config = getApplicationConfig_Unknown10()
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::GetApplicationConfig Unsupported applicationID: %v\n", applicationID)
	}

	configNative := make(types.List[types.UInt32], 0, len(config))
	for i := range config {
		configNative = append(configNative, types.NewUInt32(config[i]))
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	configNative.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodGetApplicationConfig
	rmcResponse.CallID = callID

	return rmcResponse, nil
}

func getApplicationConfig_PlayerConfig() []uint32 {
	// * This seems to be per-user configuration
	// * settings, based on the fact that the
	// * number of max uploads a user can do is
	// * sent here. No idea what anything else
	// * means
	return []uint32{
		0x00000001, 0x00000032, 0x00000096, 0x0000012c, 0x000001f4,
		0x00000320, 0x00000514, 0x000007d0, 0x00000bb8, 0x00001388,
		MAX_COURSE_UPLOADS, 0x00000014, 0x0000001e, 0x00000028, 0x00000032,
		0x0000003c, 0x00000046, 0x00000050, 0x0000005a, 0x00000064,
		0x00000023, 0x0000004b, 0x00000023, 0x0000004b, 0x00000032,
		0x00000000, 0x00000003, 0x00000003, 0x00000064, 0x00000006,
		0x00000001, 0x00000060, 0x00000005, 0x00000060, 0x00000000,
		0x000007e4, 0x00000001, 0x00000001, 0x0000000c, 0x00000000,
	}
}

func getApplicationConfig_OfficialMakers() []uint32 {
	// * Used as the PIDs for the "Official" makers in the "MAKERS" section
	return []uint32{
		2,          // * Not a real user PID, this translates to the internal Quazal Rendez-Vous user used by NEX
		1770179696, // * "official_player0" on NN, need to make PN versions
		1770179664, // * "official_player1" on NN, need to make PN versions
		1770179640, // * "official_player2" on NN, need to make PN versions
		1770180827, // * "official_player3" on NN, need to make PN versions
		1770180777, // * "official_player4" on NN, need to make PN versions
		1770180745, // * "official_player5" on NN, need to make PN versions
		1770177625, // * "official_player6" on NN, need to make PN versions
		1770177590, // * "official_player7" on NN, need to make PN versions
	}
}

func getApplicationConfig_Unknown2() []uint32 {
	// * I have no idea what this is
	// * Just replaying data sent from the real server
	return []uint32{0x000007df, 0x0000000c, 0x00000016, 0x00000005, 0x00000000}
}

func getApplicationConfig_Unknown10() []uint32 {
	// * I have no idea what this is
	// * Just replaying data sent from the real server
	// * Only seen on the 3DS
	return []uint32{35, 75, 96, 40, 5, 6}
}
