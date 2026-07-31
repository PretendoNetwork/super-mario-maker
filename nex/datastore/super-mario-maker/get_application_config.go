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
	case 0: // * Global config?
		config = getApplicationConfig_GlobalConfig()
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

func getApplicationConfig_GlobalConfig() []uint32 {
	// * This seems to be global configuration settings
	return []uint32{
		1,                  // * Number of stars for the 1st medal
		50,                 // * Number of stars for the 2nd medal
		150,                // * Number of stars for the 3rd medal
		300,                // * Number of stars for the 4th medal
		500,                // * Number of stars for the 5th medal
		800,                // * Number of stars for the 6th medal
		1300,               // * Number of stars for the 7th medal
		2000,               // * Number of stars for the 8th medal
		3000,               // * Number of stars for the 9th medal
		5000,               // * Number of stars for the 10th medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload when you have 0 or 1 medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 2nd medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 3rd medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 4th medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 5th medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 6th medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 7th medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 8th medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 9th medal
		MAX_COURSE_UPLOADS, // * Number of courses you can upload after getting the 10th medal
		35, 75, 35, 75, 50, // * Unknown
		0, 3, 3, 100, 6,    // * Unknown
		1, 96, 5, 96, 0,    // * Unknown
		2020, 1, 1, 12, 0,  // * Looks like a date? 2020-01-01 12:00?
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
	// * I have no idea what this is, looks like a date?
	// * This was when the SMM bookmark was released, so maybe it controls accessibility to it?
	// * Just replaying data sent from the real server
	return []uint32{2015, 12, 22, 5, 0} // * 2015-12-22 5:00
}

func getApplicationConfig_Unknown10() []uint32 {
	// * I have no idea what this is
	// * Just replaying data sent from the real server
	// * Only seen on the 3DS
	return []uint32{35, 75, 96, 40, 5, 6}
}
