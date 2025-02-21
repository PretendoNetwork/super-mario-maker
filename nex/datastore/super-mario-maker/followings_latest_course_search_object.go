package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func FollowingsLatestCourseSearchObject(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreSearchParam, extraData types.List[types.String]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	pRankingResults := types.NewList[datastore_super_mario_maker_types.DataStoreCustomRankingResult]()

	// * This seems to ONLY be used to get rankings for course objects
	// * uploaded by the users in param.OwnerIDs? If param.ResultOption
	// * contains the flag 0x20 ("return scores"), then the real server
	// * does some kind of check over extraData. It's unknown what this
	// * check is, so it's not done here. All other data in param seems
	// * to be unused here.
	for _, pid := range param.OwnerIDs {
		courseObjectIDs, nexError := datastore_smm_db.GetUserCourseObjectIDs(pid)
		if nexError != nil {
			return nil, nexError
		}

		// * This method seems to always use slot 0?
		results := datastore_smm_db.GetCustomRankingsByDataIDs(types.NewUInt32(0), courseObjectIDs)

		for _, rankingResult := range results {
			// * This is kind of backwards.
			// * The database pulls this data
			// * by default, so it can be done
			// * in a single query. So instead
			// * of checking if a flag *IS*
			// * set, and conditionally *ADDING*
			// * the fields, we check if a flag
			// * is *NOT* set and conditionally
			// * *REMOVE* the field
			if param.ResultOption&0x1 == 0 {
				rankingResult.MetaInfo.Tags = types.NewList[types.String]()
			}

			if param.ResultOption&0x2 == 0 {
				rankingResult.MetaInfo.Ratings = types.NewList[datastore_types.DataStoreRatingInfoWithSlot]()
			}

			if param.ResultOption&0x4 == 0 {
				rankingResult.MetaInfo.MetaBinary = types.NewQBuffer(nil)
			}

			// TODO - If this flag is set, extraData is checked somehow
			if param.ResultOption&0x20 == 0 {
				rankingResult.Score = 0
			}

			pRankingResults = append(pRankingResults, rankingResult)
		}
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pRankingResults.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodFollowingsLatestCourseSearchObject
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
