package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func CTRPickUpCourseSearchObject(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreSearchParam, extraData types.List[types.String]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// * This method is only used by the 3DS version
	// * of Super Mario Maker and is functionally
	// * identical to the WiiU versions
	// * DataStoreSMM::RecommendedCourseSearchObject
	// * method. The 3DS version LIKELY uses this as
	// * a way to filter out courses which use the
	// * Mystery Mushroom powerup, since it's not
	// * officially supported by the 3DS version. That
	// * said, the Mystery Mushroom IS still in the
	// * game and can be used somewhat normally, so
	// * no filtering is done here to prevent that.
	// * I'm not even sure how we would detect that

	// TODO - Research extraData
	// TODO - Use the offet? Real client never uses it, but might be nice for completeness sake?
	pRankingResults, nexError := datastore_smm_db.GetRandomCoursesWithLimit(int(param.ResultRange.Length))
	if nexError != nil {
		return nil, nexError
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pRankingResults.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodCTRPickUpCourseSearchObject
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
