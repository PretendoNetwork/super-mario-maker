package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func RecommendedCourseSearchObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreSearchParam, extraData []string) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	// * This method is used in 100 Mario and Course World
	// *
	// * extraData seems to be a set of filters defining a
	// * range for a courses success rate
	// *
	// * Course World (All)          ["",  "",   "",    "0", "0"]
	// * Course World (Easy)         ["1", "0",  "34",  "0", "0"]
	// * Course World (Normal)       ["1", "35", "74",  "0", "0"]
	// * Course World (Expert)       ["1", "75", "95",  "0", "0"]
	// * Course World (Super Expert) ["1", "96", "100", "0", "0"]
	// *
	// * Indexes 1 and 2 seem to be a min and max for the *failure*
	// * rate of the courses. This is not taken into account yet,
	// * as the SQL query for this would need to be rather complex.
	// * The last 2 values always seem to be 0, and the first seems
	// * to always be 1 besides filtering for "All"
	// *
	// ! All requests are treated as filtering for "All" right now
	// TODO - Use these ranges to properly filter by difficulty

	// TODO - Use the offet? Real client never uses it, but might be nice for completeness sake?
	pRankingResults, errCode := datastore_smm_db.GetRandomCoursesWithLimit(int(param.ResultRange.Length))
	if errCode != 0 {
		return errCode
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodRecommendedCourseSearchObject, rmcResponseBody)

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
