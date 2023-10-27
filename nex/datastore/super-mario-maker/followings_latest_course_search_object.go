package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func FollowingsLatestCourseSearchObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreSearchParam, extraData []string) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	pRankingResults := make([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, 0)

	// * This seems to ONLY be used to get rankings for course objects
	// * uploaded by the users in param.OwnerIDs? If param.ResultOption
	// * contains the flag 0x20 ("return scores"), then the real server
	// * does some kind of check over extraData. It's unknown what this
	// * check is, so it's not done here. All other data in param seems
	// * to be unused here.
	for _, pid := range param.OwnerIDs {
		courseObjectIDs, errCode := datastore_smm_db.GetUserCourseObjectIDs(pid)
		if errCode != 0 {
			return errCode
		}

		// * This method seems to always use slot 0?
		results := datastore_smm_db.GetCustomRankingsByDataIDs(0, courseObjectIDs)

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
				rankingResult.MetaInfo.Tags = make([]string, 0)
			}

			if param.ResultOption&0x2 == 0 {
				rankingResult.MetaInfo.Ratings = make([]*datastore_types.DataStoreRatingInfoWithSlot, 0)
			}

			if param.ResultOption&0x4 == 0 {
				rankingResult.MetaInfo.MetaBinary = make([]byte, 0)
			}

			// TODO - If this flag is set, extraData is checked somehow
			if param.ResultOption&0x20 == 0 {
				rankingResult.Score = 0
			}

			pRankingResults = append(pRankingResults, rankingResult)
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodFollowingsLatestCourseSearchObject, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
