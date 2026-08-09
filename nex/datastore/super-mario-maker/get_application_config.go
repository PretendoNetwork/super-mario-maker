package nex_datastore_super_mario_maker

import (
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
	"github.com/PretendoNetwork/super-mario-maker/nex/datastore/super-mario-maker/constants"
)

func GetApplicationConfig(err error, packet nex.PacketInterface, callID uint32, applicationID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	var config []int32

	switch applicationID {
	case 0: // * Gameplay config (Wii U)
		config = getApplicationConfig_GameplayConfig()
	case 1: // * PIDs of the "Official" makers in the "MAKERS" section
		config = getApplicationConfig_OfficialMakers()
	case 2: // * SMM bookmark
		config = getApplicationConfig_Bookmark()
	case 10: // * Gameplay config (3DS)
		config = getApplicationConfig_GameplayConfig3DS()
	default:
		fmt.Printf("[Warning] DataStoreSMMProtocol::GetApplicationConfig Unsupported applicationID: %v\n", applicationID)
	}

	configNative := make(types.List[types.Int32], 0, len(config))
	for i := range config {
		configNative = append(configNative, types.Int32(config[i]))
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	configNative.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodGetApplicationConfig
	rmcResponse.CallID = callID

	return rmcResponse, nil
}

func getApplicationConfig_GameplayConfig() []int32 {
	// * This seems to be gameplay configuration settings
	return []int32{
		constants.STARS_1ST_MEDAL,
		constants.STARS_2ND_MEDAL,
		constants.STARS_3RD_MEDAL,
		constants.STARS_4TH_MEDAL,
		constants.STARS_5TH_MEDAL,
		constants.STARS_6TH_MEDAL,
		constants.STARS_7TH_MEDAL,
		constants.STARS_8TH_MEDAL,
		constants.STARS_9TH_MEDAL,
		constants.STARS_10TH_MEDAL,

		constants.MAX_COURSE_UPLOADS_0TH_1ST_MEDAL,
		constants.MAX_COURSE_UPLOADS_2ND_MEDAL,
		constants.MAX_COURSE_UPLOADS_3RD_MEDAL,
		constants.MAX_COURSE_UPLOADS_4TH_MEDAL,
		constants.MAX_COURSE_UPLOADS_5TH_MEDAL,
		constants.MAX_COURSE_UPLOADS_6TH_MEDAL,
		constants.MAX_COURSE_UPLOADS_7TH_MEDAL,
		constants.MAX_COURSE_UPLOADS_8TH_MEDAL,
		constants.MAX_COURSE_UPLOADS_9TH_MEDAL,
		constants.MAX_COURSE_UPLOADS_10TH_MEDAL,

		// * These values are most likely settings about requesting courses, but it's not certain what they individually mean (mostly)
		constants.COURSE_WORLD_NORMAL_FAILURE_RATE_MIN,
		constants.COURSE_WORLD_EXPERT_FAILURE_RATE_MIN,
		constants.COURSE_WORLD_NORMAL_FAILURE_RATE_MIN,
		constants.COURSE_WORLD_EXPERT_FAILURE_RATE_MIN,
		50,
		0,
		3,
		3,
		constants.COURSE_DOWNLOAD_WIIU, // * Maybe?
		6, // * Might be related with the extraData?
		1,
		constants.COURSE_WORLD_SUPER_EXPERT_FAILURE_RATE_MIN,
		5,
		constants.COURSE_WORLD_SUPER_EXPERT_FAILURE_RATE_MIN,
		0,

		// * Looks like a date, possibly when the config was last changed? 2020-01-01 12:00?
		constants.CHANGED_DATE_YEAR,
		constants.CHANGED_DATE_MONTH,
		constants.CHANGED_DATE_DAY,
		constants.CHANGED_DATE_HOUR,
		constants.CHANGED_DATE_MINUTE,
	}
}

func getApplicationConfig_OfficialMakers() []int32 {
	// * Used as the PIDs for the "Official" makers in the "MAKERS" section
	return []int32{
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

func getApplicationConfig_Bookmark() []int32 {
	// * Looks like a date?
	// * This was when the SMM bookmark was released, so maybe it controls accessibility to it?
	// * Just replaying data sent from the real server
	return []int32{
		constants.BOOKMARK_DATE_YEAR,
		constants.BOOKMARK_DATE_MONTH,
		constants.BOOKMARK_DATE_DAY,
		constants.BOOKMARK_DATE_HOUR,
		constants.BOOKMARK_DATE_MINUTE,
	} // * 2015-12-22 5:00
}

func getApplicationConfig_GameplayConfig3DS() []int32 {
	// * This seems to be gameplay configuration settings for the 3DS
	return []int32{
		constants.COURSE_WORLD_NORMAL_FAILURE_RATE_MIN,
		constants.COURSE_WORLD_EXPERT_FAILURE_RATE_MIN,
		constants.COURSE_WORLD_SUPER_EXPERT_FAILURE_RATE_MIN,
		constants.COURSE_DOWNLOAD_3DS, // * Probably
		5,  // * Unknown. Might be the resultOption value?
		6,  // * Unknown. Might be related with the extraData?
	}
}
